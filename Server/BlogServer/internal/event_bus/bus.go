package event_bus

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
)

type HandlerFunc func(ctx context.Context, evt *messaging.OutboxEvent) error

type Bus struct {
	handlers map[string][]HandlerFunc
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]HandlerFunc),
	}
}

func (b *Bus) Subscribe(eventName string, handler HandlerFunc) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *Bus) Publish(ctx context.Context, evt *messaging.OutboxEvent) []error {
	var errs []error
	if handlers, ok := b.handlers[evt.EventType]; ok {
		for _, h := range handlers {
			if err := h(ctx, evt); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
