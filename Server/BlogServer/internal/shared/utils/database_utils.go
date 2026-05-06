package utils

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

func GetExecutor(ctx context.Context, db DBTX) DBTX {
	if tx, ok := ctx.Value(database.TxKey{}).(pgx.Tx); ok {
		return tx
	}
	return db
}

func GetStringPointerFromText(text pgtype.Text) *string {
	if text.Valid {
		return &text.String
	}
	return nil
}

func GetTextTypeFromNullableString(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{
			Valid: false,
		}
	}
	return pgtype.Text{
		Valid:  true,
		String: *s,
	}
}

func GetInt32PointerFromInt4(n pgtype.Int4) *int32 {
	if n.Valid {
		return &n.Int32
	}
	return nil
}

func GetFloat64PointerFromFloat8(f pgtype.Float8) *float64 {
	if f.Valid {
		return &f.Float64
	}
	return nil
}
