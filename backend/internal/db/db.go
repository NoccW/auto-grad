package db

import (
	"auto-grad-backend/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL is empty")
	}
	pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping postgres failed: %w", err)
	}
	if err := ensureTables(pool); err != nil {
		return nil, err
	}
	return pool, nil
}

func ensureTables(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
  username TEXT NOT NULL,
  role TEXT NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  student_name TEXT,
  class TEXT,
  school TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  PRIMARY KEY (username, role)
);

CREATE TABLE IF NOT EXISTS gradings (
  id TEXT PRIMARY KEY,
  subject TEXT,
  paper_image TEXT,
  answer_image TEXT,
  description TEXT,
  status TEXT,
  score INT DEFAULT 0,
  ai_score INT DEFAULT 0,
  total_score INT DEFAULT 100,
  submit_time TIMESTAMPTZ,
  created_at TIMESTAMPTZ,
  complete_time TIMESTAMPTZ,
  feedback TEXT,
  ocr_result TEXT,
  owner_username TEXT,
  owner_role TEXT
);
`)
	return err
}
