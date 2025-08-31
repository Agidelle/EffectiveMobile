package storage

import (
	"context"
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/config"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, cfg *config.Config) *Storage {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	con, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}
	// Настройки пула соединений
	con.MaxConns = 10
	con.MinConns = 1
	con.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, con)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Storage{pool: pool}
}

func (s *Storage) CloseDB() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *Storage) Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
	subs := make([]*domain.Subscription, 0)
	query := "SELECT user_id, service_name, price, start_date, end_date FROM subscriptions"
	args := []interface{}{}
	conditions := []string{}
	argIdx := 1

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}
	if filter.ServiceName != nil {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(argIdx))
		args = append(args, *filter.ServiceName)
		argIdx++
	}
	if filter.Price != nil {
		conditions = append(conditions, "price = $"+strconv.Itoa(argIdx))
		args = append(args, *filter.Price)
		argIdx++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, "start_date = $"+strconv.Itoa(argIdx))
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, "end_date = $"+strconv.Itoa(argIdx))
		args = append(args, *filter.EndDate)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if filter.Limit != nil {
		query += " LIMIT $" + strconv.Itoa(argIdx)
		args = append(args, *filter.Limit)
		argIdx++
	}
	if filter.Offset != nil {
		query += " OFFSET $" + strconv.Itoa(argIdx)
		args = append(args, *filter.Offset)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Subscription
		err = rows.Scan(&s.UserID, &s.ServiceName, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *Storage) Create(ctx context.Context, sub *domain.Subscription) error {
	res, err := s.pool.Exec(ctx,
		"INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5)",
		sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
	if err != nil {
		slog.Error("Error inserting subscription", "error", err)
		return err
	}
	if res.RowsAffected() != 1 {
		errMsg := fmt.Sprintf("expected to affect 1 row, affected: %d", res.RowsAffected())
		slog.Error(errMsg)
		return fmt.Errorf(errMsg)
	}
	slog.Info("Subscription created successfully", "user_id", sub.UserID, "service_name", sub.ServiceName)
	return nil
}

func (s *Storage) Update(ctx context.Context, sub *domain.Subscription) error {
	query := "UPDATE subscriptions SET price = $1, start_date = $2"
	args := []interface{}{sub.Price, sub.StartDate}
	argIdx := 3

	if sub.EndDate != nil {
		query += ", end_date = $" + strconv.Itoa(argIdx)
		args = append(args, *sub.EndDate)
		argIdx++
	}

	query += " WHERE user_id = $" + strconv.Itoa(argIdx) + " AND service_name = $" + strconv.Itoa(argIdx+1)
	args = append(args, sub.UserID, sub.ServiceName)

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("Error updating subscription", "error", err)
		return err
	}

	if res.RowsAffected() == 0 {
		slog.Warn("No subscription found to update", "user_id", sub.UserID, "service_name", sub.ServiceName)
		return fmt.Errorf("no subscription found for user %s and service %s", sub.UserID, sub.ServiceName)
	}
	slog.Info("Subscription updated successfully", "user_id", sub.UserID, "service_name", sub.ServiceName)
	return nil
}

// Delete подразумевается, что пользователь отменяет подписку и не важны сроки ее действия
func (s *Storage) Delete(ctx context.Context, filter *domain.Filter) error {
	res, err := s.pool.Exec(ctx, "DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2", filter.UserID, filter.ServiceName)
	if err != nil {
		slog.Error("Error deleting subscription", "error", err)
		return err
	}
	if res.RowsAffected() == 0 {
		slog.Warn("No subscription found to delete", "user_id", filter.UserID, "service_name", filter.ServiceName)
		return fmt.Errorf("subscription not found")
	}
	slog.Info("Subscription deleted successfully", "user_id", filter.UserID, "service_name", filter.ServiceName)
	return nil
}

func (s *Storage) GetSubscriptionsForPeriod(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
	query := `
		SELECT user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)
		  AND ($3::text IS NULL OR user_id = $3)
		  AND ($4::text IS NULL OR service_name = $4)
	`
	args := []interface{}{filter.EndDate, filter.StartDate, filter.UserID, filter.ServiceName}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Error querying subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(&sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate); err != nil {
			slog.Error("Error scanning subscription", "error", err)
			return nil, err
		}
		subs = append(subs, &sub)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows", "error", err)
		return nil, err
	}

	return subs, nil
}
