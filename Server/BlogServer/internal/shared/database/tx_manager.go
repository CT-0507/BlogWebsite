package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	WithVoidTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &txManager{
		pool: pool,
	}
}

type TxKey struct{}

func (tm *txManager) WithVoidTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, TxKey{}, tx)

	if err := fn(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func WithTx[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(ctx context.Context) (T, error),
) (T, error) {

	var zero T

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return zero, err
	}

	ctx = context.WithValue(ctx, TxKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	result, err := fn(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return zero, err
	}

	return result, nil
}
