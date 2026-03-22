package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a new PostgreSQL connection pool.
func NewPostgresPool(ctx context.Context, dsn string, maxConns int, maxIdleTime time.Duration) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	config.MaxConns = int32(maxConns)
	config.MaxConnIdleTime = maxIdleTime
	config.MinConns = 2
	// Disable prepared statement cache for PgBouncer (Supabase) compatibility
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// SetSchoolContext sets the current school_id for RLS (Row Level Security).
// Must be called at the start of each request for multi-tenant isolation.
func SetSchoolContext(ctx context.Context, pool *pgxpool.Pool, schoolID string) error {
	_, err := pool.Exec(ctx, "SELECT set_config('app.current_school_id', $1, true)", schoolID)
	if err != nil {
		return fmt.Errorf("set school context: %w", err)
	}
	return nil
}
