package event

import "context"

type HandlerFunc func(ctx context.Context, payload []byte) error

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

func (b *Bus) Publish(ctx context.Context, eventName string, payload []byte) error {
	if handlers, ok := b.handlers[eventName]; ok {
		for _, h := range handlers {
			if err := h(ctx, payload); err != nil {
				return err
			}
		}
	}
	return nil
}
