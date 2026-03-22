package mocks_test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockTxManager struct {
	mock.Mock
}

func (tm *MockTxManager) WithVoidTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := fn(ctx); err != nil {
		return err
	}
	return nil
}

func WithTx[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var zero T

	result, err := fn(ctx)
	if err != nil {
		return zero, err
	}

	return result, nil
}
