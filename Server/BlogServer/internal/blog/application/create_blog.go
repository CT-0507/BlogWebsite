package application

import (
	"context"
	"fmt"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/google/uuid"
)

type CreateBlogUseCases struct {
	txManager   database.TxManager
	repo        domain.BlogRepository
	userService user.UserService
	outboxRepo  outbox.OutboxRepository
}

func NewCreateBlogUseCases(txManager database.TxManager, repo domain.BlogRepository, userService user.UserService, outboxRepo outbox.OutboxRepository) *CreateBlogUseCases {
	return &CreateBlogUseCases{
		txManager:   txManager,
		repo:        repo,
		userService: userService,
		outboxRepo:  outboxRepo,
	}
}

// Save a box to database and Create an Event to outbox_events table
func (s *CreateBlogUseCases) CreateWithOutBox(c context.Context, blog *domain.Blog, userID string) error {

	authorID, err := s.repo.VerifyAuthorIDByUserID(c, userID)
	if err != nil {
		return err
	}
	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		blog.AuthorID = authorID

		insertedBlog, err := s.repo.Create(c, blog)
		if err != nil {
			return err
		}

		// Save event to outbox table
		event := &BlogCreatedEvent{
			BlogID:    insertedBlog.BlogID,
			AuthorID:  authorID,
			BlogTitle: insertedBlog.Title,
		}

		err = s.outboxRepo.Insert(c, &messaging.OutboxEvent{
			Payload: event,
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
func (s *CreateBlogUseCases) OnBlogPosted(c context.Context, payload BlogCreatedEvent) error {
	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		content := fmt.Sprintf("A blog with title %s has just been created", payload.BlogTitle)
		not, err := s.userService.CreateNotification(c, content, uuid.MustParse(config.ADMIN_ID), uuid.MustParse(config.SYSTEM_ID))
		if err != nil {
			return err
		}

		newPayload := map[string]any{
			"blogID":    payload.BlogID,
			"blogTitle": payload.BlogTitle,
			"userID":    not.UserID,
		}

		err = s.outboxRepo.Insert(c, &messaging.OutboxEvent{
			EventType: "blogCreated",
			Payload:   newPayload,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
