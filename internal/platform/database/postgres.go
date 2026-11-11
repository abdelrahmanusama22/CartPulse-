// Package database owns connection pools; business transactions belong to services.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		// Do not return the connection URL: it can contain credentials.
		return nil, fmt.Errorf("invalid DATABASE_URL")
	}
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	cfg.MaxConns = 10
	cfg.MinConns = 0
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize PostgreSQL connection pool")
	}
	// Connections are lazy. The readiness endpoint performs an actual Ping.
	return pool, nil
}
