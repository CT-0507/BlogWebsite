package notification

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/sse"
)

type NotificationPublisher interface {
	PublishCache(topic string, event string, data *sse.Cache)
	PublishEvent(topic string, event string, data *interface{})
}

type NotificationService struct {
	publisher NotificationPublisher
}

func NewNotificationService(publisher NotificationPublisher) *NotificationService {
	return &NotificationService{publisher: publisher}
}

// type NotificationCreatedEvent struct {
// 	NotID   string
// 	Content string
// }

func (s *NotificationService) PublishNotification(c context.Context, evt *messaging.OutboxEvent) error {

	var event interface{}
	if err := json.Unmarshal(evt.Payload, &event); err != nil {
		return err
	}
	s.publisher.PublishCache("blog_created_admin", "blog", &sse.Cache{
		QueryKey: []string{"notifications"},
		Op:       "append",
		Data:     event,
	})

	return nil
}
