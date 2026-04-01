package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type CreateBlogUseCases struct {
	txManager  database.TxManager
	repo       domain.BlogRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewCreateBlogUseCases(
	txManager database.TxManager,
	repo domain.BlogRepository,
	outboxRepo outboxrepo.OutboxRepository,
) *CreateBlogUseCases {
	return &CreateBlogUseCases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
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

func (s *CreateBlogUseCases) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return s.repo.VerifyAuthorIDByUserID(c, userID)
}
