package database

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

type QueryTracer struct{}

func shouldSkip(sql string) bool {
	// simple check (case-insensitive is safer in real code)
	s := strings.ToLower(strings.TrimSpace(sql))

	// skip your table
	if strings.Contains(s, "outbox.outbox_events") {
		return true
	}

	// skip transaction noise
	switch s {
	case "begin", "commit", "rollback":
		return true
	}

	return false
}

func (t *QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if shouldSkip(data.SQL) {
		return ctx
	}
	log.Printf("START SQL: %s, args=%v\n", data.SQL, data.Args)
	return ctx
}

func (t *QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if shouldSkip(data.CommandTag.String()) {
		return
	}

	// log.Printf("END SQL: %s, err=%v", data.CommandTag, data.Err)
}
