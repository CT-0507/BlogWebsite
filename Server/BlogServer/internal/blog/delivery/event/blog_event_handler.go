package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
)

type EventHandler struct {
	txManager  database.TxManager
	repo       domain.BlogRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewEventHandler(txManager database.TxManager, repo domain.BlogRepository, outboxRepo outboxrepo.OutboxRepository) *EventHandler {
	return &EventHandler{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (e *EventHandler) OnAuthorCreated(c context.Context, evt *messaging.OutboxEvent) error {

	var stepPayload contracts.CreateBlogAuthorCachePayload
	_ = json.Unmarshal(evt.Payload, &stepPayload)

	err := e.txManager.WithVoidTx(c, func(ctx context.Context) error {

		err := e.repo.CreateUserIDAuthorProfileIDCacheRecord(c, stepPayload.UserID, stepPayload.AuthorID, stepPayload.Slug, stepPayload.DisplayName)
		if err != nil {
			return err
		}

		eventPayload := map[string]any{}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.CreateBlogAuthorCacheSuccessContext{
			UserID:   stepPayload.UserID,
			AuthorID: stepPayload.AuthorID,
		}

		context, _ := json.Marshal(eventContext)

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "CreateBlogAuthorCache.Success",
			Payload:   payload,
			Context:   &context,
		})
	})

	// Signal failed saga step
	if err != nil {

		ctx := context.WithValue(c, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "CreateBlogAuthorCache.Failed",
			Payload:   b,
			Error:     evt.Error,
		})
		if err1 != nil {
			return err1
		}
		return err
	}
	return nil
}

func (e *EventHandler) OnDeleteBlogAuthorCache(c context.Context, evt *messaging.OutboxEvent) error {

	var outboxPayload contracts.DeleteAuthorPayload
	err := json.Unmarshal(evt.Payload, &outboxPayload)
	if err != nil {
		return err
	}

	err = e.txManager.WithVoidTx(c, func(ctx context.Context) error {

		err := e.repo.UpdateBlogStatusForDeletedAuthor(c, outboxPayload.AuthorID)
		if err != nil {
			return err
		}

		err = e.repo.MarkAuthorCacheAsDeleted(c, outboxPayload.AuthorID)
		if err != nil {
			return err
		}

		eventPayload := map[string]any{}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.DeleteBlogAuthorCacheContext{
			AuthorID: outboxPayload.AuthorID,
		}

		context, _ := json.Marshal(eventContext)

		return e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "DeleteBlogAuthorCache.Success",
			Payload:   payload,
			Context:   &context,
		})
	})

	// Signal failed saga step
	if err != nil {

		ctx := context.WithValue(c, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "DeleteBlogAuthorCache.Failed",
			Payload:   b,
			Error:     evt.Error,
		})
		if err1 != nil {
			return err1
		}
		return err
	}
	return nil
}

// func (e *EventHandler) OnAuthorDeleted(c context.Context, evt *messaging.OutboxEvent) error {
// 	return e.txManager.WithVoidTx(c, func(ctx context.Context) error {

// 		var event domain.AuthorDeletedEvent
// 		if err := json.Unmarshal(evt.Payload, &event); err != nil {
// 			return err
// 		}
// 		return e.repo.UpdateBlogStatusForDeletedAuthor(c, event.AuthorID)
// 	})
// }

func (e *EventHandler) OnAuthorHardDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return e.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorDeletedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}
		err := e.repo.DeleteAuthorCache(c, event.AuthorID)
		if err != nil {
			return err
		}
		return e.repo.DeleteAuthorHardDeletedBlogs(c, event.AuthorID)
	})
}

func (e *EventHandler) CreateBlog(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	var payload contracts.BlogCreatedSagaPayload
	err := json.Unmarshal(evt.Payload, &payload)
	if err != nil {
		return err
	}

	err = e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		insertedBlog, err := e.repo.Create(ctx, &domain.Blog{
			AuthorID: payload.AuthorID,
			Title:    payload.Title,
			Content:  payload.Content,
			Status:   payload.Status,
			URLSlug:  payload.UrlSlug,
		})
		if err != nil {
			return err
		}

		truncatedTitle := utils.Truncate(payload.Title, 20, true)
		truncatedContent := utils.Truncate(payload.Content, 50, true)

		// Success
		context := &map[string]interface{}{
			"BlogID": insertedBlog.BlogID,
		}

		payload := &map[string]interface{}{
			"BlogID":           insertedBlog.BlogID,
			"AuthorID":         payload.AuthorID,
			"UserID":           payload.UserID,
			"TruncatedTitle":   truncatedTitle,
			"TruncatedContent": truncatedContent,
			"UrlSlug":          payload.UrlSlug,
		}

		payloadMarshal, _ := json.Marshal(payload)
		contextMarshal, _ := json.Marshal(context)
		// Save event to outbox table
		// event := &contracts.BlogCreatedEvent{
		// 	BlogID:    insertedBlog.BlogID,
		// 	AuthorID:  authorID,
		// 	BlogTitle: insertedBlog.Title,
		// }

		// Proceed next step
		err = e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "InceaseAuthorBlogCount",
			Payload:   payloadMarshal,
			Context:   &contextMarshal,
		})
		if err != nil {
			return err
		}

		// Create notification asynchronously
		blogCreatedEvt := &contracts.BlogCreatedEvent{
			BlogID:    insertedBlog.BlogID,
			BlogTitle: truncatedTitle,
			AuthorID:  insertedBlog.AuthorID,
		}

		blogCreatedEvtPayload, _ := json.Marshal(blogCreatedEvt)
		err = e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: "blog.created",
			Payload:   blogCreatedEvtPayload,
		})

		if err != nil {
			return err
		}

		return nil
	})

	// Signal failed saga step
	if err != nil {

		ctx = context.WithValue(ctx, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "CreateBlog.Failed",
			Payload:   b,
			Error:     evt.Error,
		})
		if err1 != nil {
			log.Println(err1)
			return err1
		}
		return err
	}
	return nil
}

// Handle blog posted event for event bus
func (s *EventHandler) OnBlogPosted(c context.Context, evt *messaging.OutboxEvent) error {
	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		// var event contracts.BlogCreatedEvent
		json.Unmarshal(evt.Payload, &evt)

		// content := fmt.Sprintf("A blog with title %s has just been created", event.BlogTitle)
		// not, err := s.userService.CreateNotification(c, content, uuid.MustParse(config.ADMIN_ID), uuid.MustParse(config.SYSTEM_ID))
		// if err != nil {
		// 	return err
		// }

		// userUUID := uuid.MustParse(not.UserID)

		// newPayload := contracts.BlogCreatedSagaPayload{
		// 	// BlogID: event.BlogID,
		// 	Title:  event.BlogTitle,
		// 	UserID: userUUID,
		// }
		// newPayloadMarshal, _ := json.Marshal(newPayload)
		// sagaID := uuid.New()

		// newContext := contracts.BlogCreatedSagaContext{
		// 	// BlogID: event.BlogID,
		// 	// Title:  event.BlogTitle,
		// 	UserID: userUUID,
		// }
		// newContextMarshal, _ := json.Marshal(newContext)
		// err = s.outboxRepo.Insert(c, &messaging.OutboxEvent{
		// 	SagaID:    &sagaID,
		// 	EventType: "blogCreated",
		// 	Payload:   newPayloadMarshal,
		// 	Context:   &newContextMarshal,
		// })
		// if err != nil {
		// 	return err
		// }
		return nil
	})
}

func (e *EventHandler) OnDeleteBlog(c context.Context, evt *messaging.OutboxEvent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var payload contracts.DeleteBlogPayload
	err := json.Unmarshal(evt.Payload, &payload)
	if err != nil {
		return err
	}

	err = e.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		_, err := e.repo.Delete(ctx, payload.BlogID, config.SYSTEM_ID)
		if err != nil {
			return err
		}

		// empty payload
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err = e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "DeleteBlog.Failed",
			Payload:   b,
		})
		if err != nil {
			return err
		}

		return nil
	})

	// Signal failed saga step
	if err != nil {

		ctx = context.WithValue(ctx, database.TxKey{}, nil)
		// Fail to create blog
		m := map[string]any{}
		b, _ := json.Marshal(m)
		err1 := e.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    evt.SagaID,
			EventType: "DeleteBlog.Failed",
			Payload:   b,
			Error:     evt.Error,
		})
		if err1 != nil {
			log.Println(err1)
			return err1
		}
		return err
	}
	return nil
}
