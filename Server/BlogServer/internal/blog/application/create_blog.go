package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/google/uuid"
)

type CreateBlogUseCases struct {
	txManager   database.TxManager
	repo        domain.BlogRepository
	userService user.UserService
	outboxRepo  outboxrepo.OutboxRepository
}

func NewCreateBlogUseCases(txManager database.TxManager, repo domain.BlogRepository, userService user.UserService, outboxRepo outboxrepo.OutboxRepository) *CreateBlogUseCases {
	return &CreateBlogUseCases{
		txManager:   txManager,
		repo:        repo,
		userService: userService,
		outboxRepo:  outboxRepo,
	}
}

// Save a box to database and Create an Event to outbox_events table
func (s *CreateBlogUseCases) CreateBlogStartSaga(c context.Context, blog *domain.Blog, userID string) error {

	authorID, err := s.repo.VerifyAuthorIDByUserID(c, userID)
	if err != nil {
		return err
	}
	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		blog.AuthorID = authorID

		userUUID := uuid.MustParse(userID)
		// Save event to outbox table
		context := &contracts.BlogCreatedSagaContext{
			AuthorID: authorID,
			UserID:   userUUID,
		}

		payload := &contracts.BlogCreatedSagaPayload{
			AuthorID: authorID,
			UserID:   userUUID,
			Title:    blog.Title,
			UrlSlug:  blog.URLSlug,
			Content:  blog.Content,
			Status:   blog.Status,
		}

		payloadMarshal, _ := json.Marshal(payload)
		contextMarshal, _ := json.Marshal(context)
		sagaID := uuid.New()
		err = s.outboxRepo.Insert(c, &messaging.OutboxEvent{
			SagaID:    &sagaID,
			EventType: "create_blog_saga",
			Payload:   payloadMarshal,
			Context:   &contextMarshal,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *CreateBlogUseCases) CreateBlog(c context.Context, evt *messaging.OutboxEvent) error {

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

func (s *CreateBlogUseCases) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return s.repo.VerifyAuthorIDByUserID(c, userID)
}

// Handle blog posted event for event bus
func (s *CreateBlogUseCases) OnBlogPosted(c context.Context, evt *messaging.OutboxEvent) error {
	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event contracts.BlogCreatedEvent
		json.Unmarshal(evt.Payload, &evt)

		content := fmt.Sprintf("A blog with title %s has just been created", event.BlogTitle)
		not, err := s.userService.CreateNotification(c, content, uuid.MustParse(config.ADMIN_ID), uuid.MustParse(config.SYSTEM_ID))
		if err != nil {
			return err
		}

		userUUID := uuid.MustParse(not.UserID)

		newPayload := contracts.BlogCreatedSagaPayload{
			// BlogID: event.BlogID,
			Title:  event.BlogTitle,
			UserID: userUUID,
		}
		newPayloadMarshal, _ := json.Marshal(newPayload)
		sagaID := uuid.New()

		newContext := contracts.BlogCreatedSagaContext{
			// BlogID: event.BlogID,
			// Title:  event.BlogTitle,
			UserID: userUUID,
		}
		newContextMarshal, _ := json.Marshal(newContext)
		err = s.outboxRepo.Insert(c, &messaging.OutboxEvent{
			SagaID:    &sagaID,
			EventType: "blogCreated",
			Payload:   newPayloadMarshal,
			Context:   &newContextMarshal,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
