package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/google/uuid"
)

type CreateBlogUseCase struct {
	txManager   *database.TxManager
	repo        domain.BlogRepository
	userService user.UserService
	outboxRepo  outbox.OutboxRepository
}

func NewCreateBlogUseCase(txManager *database.TxManager, repo domain.BlogRepository, userService user.UserService, outboxRepo outbox.OutboxRepository) *CreateBlogUseCase {
	return &CreateBlogUseCase{
		txManager:   txManager,
		repo:        repo,
		userService: userService,
		outboxRepo:  outboxRepo,
	}
}

// Save a box to database and Create an Event to outbox_events table
func (s *CreateBlogUseCase) CreateWithOutBox(c context.Context, blog *domain.Blog) error {

	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {

		insertedBlog, err := s.repo.Create(c, blog)
		if err != nil {
			return err
		}

		event := &BlogCreatedEvent{
			BlogID:    insertedBlog.BlogID,
			BlogTitle: insertedBlog.Title,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = s.outboxRepo.Insert(c, event.EventName(), payload)
		if err != nil {
			return err
		}

		return nil
	})

}

func (s *CreateBlogUseCase) OnBlogPosted(c context.Context, payload []byte) error {

	return s.txManager.WithVoidTx(c, func(ctx context.Context) error {
		var evt BlogCreatedEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}
		content := fmt.Sprintf("A blog with title %s has just been created", evt.BlogTitle)
		not, err := s.userService.CreateNotification(c, content, uuid.MustParse(config.ADMIN_ID), uuid.MustParse(config.SYSTEM_ID))
		if err != nil {
			return err
		}

		notificationPayload, err := json.Marshal(not)
		if err != nil {
			return err
		}

		err = s.outboxRepo.Insert(c, "notification.created", notificationPayload)
		if err != nil {
			return err
		}
		return nil
	})

}
