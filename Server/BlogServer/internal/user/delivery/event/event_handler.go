package event

import (
	"context"
	"encoding/json"
	"time"

	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
)

type EventHandler struct {
	txManager  database.TxManager
	repo       domain.UserRepository
	outboxRepo outboxrepo.OutboxRepository
}

func New(txManager database.TxManager, repo domain.UserRepository, outboxRepo outboxrepo.OutboxRepository) *EventHandler {
	return &EventHandler{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (e *EventHandler) OnCreateNotifications(ctx context.Context, evt *messaging.OutboxEvent) error {
	timeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var payload CreateNotificationsEventPayload
	err := json.Unmarshal(evt.Payload, &payload)
	if err != nil {
		return err
	}

	// Process in chunk of 5
	// Ignore error
	followerIds := payload.FollowerIDs
	notificationContent := map[string]interface{}{
		"UrlSlug":    payload.UrlSlug,
		"AuthorID":   payload.AuthorID,
		"AuthorName": payload.AuthorName,
		"AuthorSlug": payload.AuthorSlug,
		"Title":      payload.TruncatedTitle,
		"Content":    payload.TruncatedContent,
	}

	contentMarshal, _ := json.Marshal(notificationContent)

	for i := 0; i < len(payload.FollowerIDs); i += 5 {
		var insertItems []domain.Notification
		for j := 0; j < 5 && i+j < len(payload.FollowerIDs); j++ {
			insertItems = append(insertItems, domain.Notification{
				UserID:  followerIds[i+j],
				Content: contentMarshal,
			})
		}
		e.repo.CreateNotifications(timeCtx, insertItems)
	}

	return nil
}
