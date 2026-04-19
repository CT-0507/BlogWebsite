package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	"github.com/google/uuid"
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

func (e *EventHandler) OnDeleteUser(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var outboxPayload contracts.DeleteUserSagaPayload
	_ = json.Unmarshal(evt.Payload, &outboxPayload)

	err := e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := e.repo.MarkUserAsDeleted(ctx, outboxPayload.UserID, outboxPayload.UpdatedBy)
		if err != nil {
			return err
		}

		eventPayload := &contracts.DeleteUserPayload{
			UserID:    outboxPayload.UserID,
			UpdatedBy: outboxPayload.UpdatedBy,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.DeleteUserContext{
			UserID:         outboxPayload.UserID,
			PreviousStatus: outboxPayload.Status,
		}

		context, err := json.Marshal(eventContext)
		if err != nil {
			return err
		}

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.DeleteUserSuccess,
			Payload:   payload,
			Context:   &context,
		})
	})

	// Signal failed saga step
	if err != nil {

		ctx = context.WithValue(ctx, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.DeleteUserFailed,
			Payload:   b,
			Error:     utils.StringPtr(err.Error()),
		})
		if err1 != nil {
			return err1
		}
		return nil
	}
	return nil
}

func (e *EventHandler) OnDeleteUserCompensation(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var outboxPayload contracts.DeleteUserContext
	_ = json.Unmarshal(evt.Payload, &outboxPayload)

	err := e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		systemUUID := uuid.MustParse(config.SYSTEM_ID)
		err := e.repo.RestoreUserByID(ctx, outboxPayload.UserID, systemUUID)
		if err != nil {
			return err
		}

		eventPayload := &map[string]any{}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.DeleteUserCompensationSuccess,
			Payload:   payload,
		})
	})

	// Signal failed saga step
	if err != nil {

		ctx = context.WithValue(ctx, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.DeleteUserCompensationFailed,
			Payload:   b,
			Error:     utils.StringPtr(err.Error()),
		})
		if err1 != nil {
			return err1
		}
		return nil
	}
	return nil
}
