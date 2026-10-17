package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	startupTimeout = 30 * time.Second
	pingTimeout    = 2 * time.Second
	retryInterval  = time.Second
)

// Postgres owns the PostgreSQL connection pool used by the API.
type Postgres struct {
	pool *pgxpool.Pool
}

// Open creates a connection pool and waits for PostgreSQL to become reachable.
func Open(ctx context.Context, connectionString string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	postgres := &Postgres{pool: pool}
	if err := postgres.waitUntilReady(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("wait for PostgreSQL: %w", err)
	}

	return postgres, nil
}

// Ping reports whether PostgreSQL is reachable within the caller's context.
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// Close releases the connection pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) waitUntilReady(ctx context.Context) error {
	deadline := time.NewTimer(startupTimeout)
	defer deadline.Stop()

	for {
		pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
		err := p.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}

		retry := time.NewTimer(retryInterval)
		select {
		case <-ctx.Done():
			retry.Stop()
			return ctx.Err()
		case <-deadline.C:
			retry.Stop()
			return err
		case <-retry.C:
		}
	}
}
