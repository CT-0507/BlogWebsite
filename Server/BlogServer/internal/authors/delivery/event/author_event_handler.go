package event

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/oklog/ulid/v2"
)

type EventHandler struct {
	txManager      database.TxManager
	repo           domain.AuthorProfileRepository
	outboxRepo     outboxrepo.OutboxRepository
	storageService storage.Storage
}

func New(txManager database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outboxrepo.OutboxRepository) *EventHandler {
	return &EventHandler{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (e *EventHandler) OnAuthorFollowerCountChanged(ctx context.Context, evt *messaging.OutboxEvent) error {

	timeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return e.txManager.WithVoidTx(timeCtx, func(ctx context.Context) error {

		var event contracts.FollowCountChangedEvent
		err := json.Unmarshal(evt.Payload, &event)
		if err != nil {
			return err
		}

		err = e.repo.UpdateAuthorFollowerCount(ctx, event.AuthorID, event.IsIncrement)
		if err != nil {
			return err
		}

		event2 := &contracts.FollowCountChangedEvent{
			AuthorID:    event.AuthorID,
			UserID:      event.UserID,
			IsIncrement: event.IsIncrement,
		}

		eventPayload, _ := json.Marshal(event2)

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType:  event2.EventName(),
			Payload:    eventPayload,
			RetryCount: 1,
		})
	})
}

func (e *EventHandler) OnAuthorCreate(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var author domain.AuthorProfile
	err := json.Unmarshal(evt.Payload, &author)
	if err != nil {
		return err
	}

	if author.Avatar != nil {

		dst := strings.Replace(*author.Avatar, "/temp/", "/", 1)

		err := e.storageService.MoveFile(*author.Avatar, dst)
		if err != nil {
			return err
		}
		author.Avatar = &dst
	}

	author.AuthorID = ulid.Make().String()

	err = e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		err := e.repo.CreateAuthorProfile(ctx, &author, author.UserID, *author.CreatedBy)
		if err != nil {
			return err
		}

		eventPayload := &contracts.CreateBlogAuthorCachePayload{
			AuthorID:    author.AuthorID,
			UserID:      author.UserID,
			Slug:        author.Slug,
			DisplayName: author.DisplayName,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.CreateBlogAuthorCacheContext{
			AuthorID: author.AuthorID,
			Avatar:   author.Avatar,
		}

		context, _ := json.Marshal(eventContext)

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.CreateAuthorSuccess,
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
			EventType: flows.CreateAuthorFailed,
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

func (e *EventHandler) OnCreateAuthorCompensation(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var output contracts.CreateBlogAuthorCacheContext
	_ = json.Unmarshal(evt.Payload, &output)

	if output.Avatar != nil {

		dst := utils.SwapTemp(*output.Avatar, true)

		err := e.storageService.MoveFile(*output.Avatar, dst)
		if err != nil {
			return err
		}
		output.Avatar = &dst
	}

	err := e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := e.repo.UpdateAuthorStatus(ctx, output.AuthorID, "deleted", config.SYSTEM_ID)
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
			EventType: flows.CreateAuthorCompensationSuccess,
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
			EventType: flows.CreateAuthorCompensationFailed,
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

func (e *EventHandler) OnDeleteAuthor(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var outboxPayload contracts.DeleteAuthorKickstartPayload
	_ = json.Unmarshal(evt.Payload, &outboxPayload)

	if outboxPayload.Avatar != nil {

		dst := utils.SwapTemp(*outboxPayload.Avatar, true)

		err := e.storageService.MoveFile(*outboxPayload.Avatar, dst)
		if err != nil {
			return err
		}
		outboxPayload.Avatar = &dst
	}

	err := e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := e.repo.UpdateAuthorStatus(ctx, outboxPayload.AuthorID, "deleted", config.SYSTEM_ID)
		if err != nil {
			return err
		}

		eventPayload := &contracts.DeleteAuthorPayload{
			AuthorID: outboxPayload.AuthorID,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.DeleteAuthorContext{
			AuthorID:       outboxPayload.AuthorID,
			PreviousStatus: outboxPayload.Status,
			Avatar:         outboxPayload.Avatar,
		}

		context, err := json.Marshal(eventContext)
		if err != nil {
			return err
		}

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.DeleteAuthorSuccess,
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
			EventType: flows.DeleteAuthorFailed,
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

func (e *EventHandler) OnDeleteAuthorCompensation(c context.Context, evt *messaging.OutboxEvent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var output contracts.DeleteAuthorContext
	_ = json.Unmarshal(evt.Payload, &output)

	if output.Avatar != nil {

		dst := utils.SwapTemp(*output.Avatar, false)

		err := e.storageService.MoveFile(*output.Avatar, dst)
		if err != nil {
			return err
		}
		output.Avatar = &dst
	}

	err := e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := e.repo.UpdateAuthorStatus(ctx, output.AuthorID, output.PreviousStatus, config.SYSTEM_ID)
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
			EventType: flows.DeleteAuthorCompensationSuccess,
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
			EventType: flows.DeleteAuthorCompensationFailed,
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

func (e *EventHandler) OnBlogCreated(ctx context.Context, evt *messaging.OutboxEvent) error {
	timeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := e.txManager.WithVoidTx(timeCtx, func(ctx context.Context) error {

		var payload contracts.BlogCountChangedEvent
		err := json.Unmarshal(evt.Payload, &payload)
		if err != nil {
			return err
		}

		err = e.repo.UpdateAuthorBlogCount(ctx, payload.AuthorID, true)
		if err != nil {
			return err
		}

		author, err := e.repo.GetAuthorProfileByID(ctx, payload.AuthorID, "active", "check_null")
		if err != nil {
			return err
		}

		context := &map[string]interface{}{
			"AuthorID": payload.AuthorID,
		}

		newPayload, err := utils.StructToMap(payload)
		if err != nil {
			return err
		}

		followerIds, err := e.repo.GetAuthorFollowersByID(ctx, payload.AuthorID)
		if err != nil {
			return err
		}

		newPayload["FollowerIds"] = followerIds
		newPayload["AuthorName"] = author.DisplayName
		newPayload["AuthorSlug"] = author.Slug

		payloadMarshal, _ := json.Marshal(newPayload)
		contextMarshal, _ := json.Marshal(context)

		err = e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: "CreateNotifications",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.InceaseAuthorBlogCountSuccess,
			Payload:   evt.Payload,
			Context:   &contextMarshal,
		})
	})

	if err != nil {
		ctx = context.WithValue(ctx, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: flows.InceaseAuthorBlogCountFailed,
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
