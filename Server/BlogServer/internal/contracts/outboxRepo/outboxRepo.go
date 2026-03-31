package outboxrepo

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
)

type OutboxRepository interface {
	Insert(ctx context.Context, evt *messaging.OutboxEvent) error
}
