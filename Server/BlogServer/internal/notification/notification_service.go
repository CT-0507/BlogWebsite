package notification

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/sse"
)

type NotificationService interface {
	PublishNotification(c context.Context, payload []byte) error
}

type notificationService struct {
	broker *sse.Broker
}

func NewNotificationService(broker *sse.Broker) NotificationService {
	return &notificationService{broker: broker}
}

type NotificationCreatedEvent struct {
	NotID   string
	Content string
}

func (s *notificationService) PublishNotification(c context.Context, payload []byte) error {
	var evt NotificationCreatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	s.broker.PublishCache("blog_created_admin", "blog", sse.Cache{
		QueryKey: []string{"notifications"},
		Op:       "append",
		Data:     evt,
	})

	return nil
}
