package storage

import (
	"context"
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/config"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
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
	return nil, nil
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
		slog.Error("expected to affect 1 row, affected:", res.RowsAffected())
		return err
	}
	slog.Info("Subscription created successfully", "user_id", sub.UserID, "service_name", sub.ServiceName)
	return nil
}

func (s *Storage) Update(ctx context.Context, sub *domain.Subscription) error {
	query := "UPDATE subscriptions SET price = $1, start_date = $2"
	args := []interface{}{sub.Price, sub.StartDate}

	if sub.EndDate != nil {
		query += ", end_date = $3"
		args = append(args, sub.EndDate)
	}

	query += " WHERE user_id = $4 AND service_name = $5"
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
func (s *Storage) Delete(ctx context.Context, sub *domain.Subscription) error {
	res, err := s.pool.Exec(ctx, "DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2", sub.UserID, sub.ServiceName)
	if err != nil {
		slog.Error("Error deleting subscription", "error", err)
		return err
	}
	if res.RowsAffected() == 0 {
		slog.Warn("No subscription found to delete", "user_id", sub.UserID, "service_name", sub.ServiceName)
		return fmt.Errorf("no subscription found for user %s and service %s", sub.UserID, sub.ServiceName)
	}
	slog.Info("Subscription deleted successfully", "user_id", sub.UserID, "service_name", sub.ServiceName)
	return nil
}
