package application

import "context"

type OutboxRepository interface {
	Insert(ctx context.Context, topic string, payload []byte) error
}
