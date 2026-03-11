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
	Message string
}

type NotificationMessage struct {
	Message string
}

func (s *notificationService) PublishNotification(c context.Context, payload []byte) error {
	var evt NotificationCreatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	s.broker.Publish("blog_created_admin", "blog", &NotificationMessage{
		Message: evt.Message,
	})

	return nil
}
