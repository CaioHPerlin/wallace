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

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS devices (
		id             text PRIMARY KEY,
		name           text NOT NULL DEFAULT '',
		led_enabled    boolean NOT NULL DEFAULT true,
		buzzer_enabled boolean NOT NULL DEFAULT true,
		last_heartbeat bigint,
		created_at     timestamptz NOT NULL DEFAULT now(),
		updated_at     timestamptz NOT NULL DEFAULT now()
	)`,
	`CREATE TABLE IF NOT EXISTS detections (
		id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		device_id text NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
		ts        bigint NOT NULL
	)`,
}

type Postgres struct {
	pool *pgxpool.Pool
}

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
	if err := postgres.ensureSchema(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ensure database schema: %w", err)
	}

	return postgres, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) ensureSchema(ctx context.Context) error {
	for _, statement := range schemaStatements {
		if _, err := p.pool.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
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
