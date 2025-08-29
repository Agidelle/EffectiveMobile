package storage

import (
	"context"
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/config"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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
		fmt.Println("Error inserting subscription:", err)
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", res.RowsAffected())
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, sub *domain.Subscription) error {
	return nil
}

func (s *Storage) Delete(ctx context.Context, id *int) error {
	return nil
}
