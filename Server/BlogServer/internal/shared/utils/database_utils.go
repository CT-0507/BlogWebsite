package utils

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func GetExecutor(ctx context.Context, db DBTX) DBTX {
	if tx, ok := ctx.Value(database.TxKey{}).(pgx.Tx); ok {
		return tx
	}
	return db
}
