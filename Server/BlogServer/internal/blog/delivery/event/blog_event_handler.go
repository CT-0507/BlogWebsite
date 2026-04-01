package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
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

func (u *EventHandler) OnAuthorCreated(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorCreatedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}

		return u.repo.CreateUserIDAuthorProfileIDCacheRecord(c, event.UserID, event.AuthorID, event.Slug, event.DisplayName)
	})
}

func (u *EventHandler) OnAuthorDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorDeletedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}
		return u.repo.UpdateBlogStatusForDeletedAuthor(c, event.AuthorID)
	})
}

func (u *EventHandler) OnAuthorHardDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorDeletedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}
		err := u.repo.DeleteAuthorCache(c, event.AuthorID)
		if err != nil {
			return err
		}
		return u.repo.DeleteAuthorHardDeletedBlogs(c, event.AuthorID)
	})
}

func (s *EventHandler) CreateBlog(c context.Context, evt *messaging.OutboxEvent) error {

	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	var payload contracts.BlogCreatedSagaPayload
	err := json.Unmarshal(evt.Payload, &payload)
	if err != nil {
		return err
	}

	return s.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		insertedBlog, err := s.repo.Create(ctx, &domain.Blog{
			AuthorID: payload.AuthorID,
			Title:    payload.Title,
			Content:  payload.Content,
			Status:   payload.Status,
			URLSlug:  payload.UrlSlug,
		})
		if err != nil {
			// Fail to create blog
			err = s.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
				SagaID:    evt.SagaID,
				EventType: "CreateBlog.Failed",
				Error:     evt.Error,
			})
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
		err = s.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
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
		err = s.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: "blog.created",
			Payload:   blogCreatedEvtPayload,
		})

		if err != nil {
			return err
		}

		return nil
	})
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
