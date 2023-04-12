package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	config Config
	*pgxpool.Pool
}

func NewDB(cfg Config) *DB {
	return &DB{config: cfg}
}

func OpenDB(ctx context.Context, cfg Config) (*DB, error) {
	db := NewDB(cfg)
	if err := db.Open(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Open(ctx context.Context) error {
	uri, err := db.config.URI()
	if err != nil {
		return err
	}

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	db.Pool = pool

	return nil
}

func (db *DB) Close() { db.Pool.Close() }
