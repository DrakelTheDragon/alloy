package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type NamedArgs = pgx.NamedArgs

type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Tx interface {
	Conn
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func Exec(ctx context.Context, conn Conn, sql string, args ...any) error {
	_, err := conn.Exec(ctx, sql, args...)
	return err
}

func Query[T any](ctx context.Context, conn Conn, sql string, args ...any) ([]T, error) {
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func QueryRow[T any](ctx context.Context, conn Conn, sql string, args ...any) (T, error) {
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return *new(T), err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}
