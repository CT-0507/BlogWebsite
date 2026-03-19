package application

import (
	"context"
	"encoding/json"
	"log"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/oklog/ulid/v2"
)

type AuthorProfileUsecases struct {
	txManager  *database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outbox.OutboxRepository
}

func NewAuthorProfileUsecases(txManager *database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outbox.OutboxRepository) *AuthorProfileUsecases {
	return &AuthorProfileUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *AuthorProfileUsecases) CreateAuthor(ctx context.Context, author *domain.AuthorProfile, userID string, createdBy string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		author.AuthorID = ulid.Make().String()
		err := u.repo.CreateAuthorProfile(ctx, author, userID, createdBy)
		if err != nil {
			log.Println(err)
			return &domain.ErrFailedToCreateAuthorProfile{}
		}

		event := &domain.AuthorCreatedEvent{
			AuthorID: author.AuthorID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *AuthorProfileUsecases) GetAuthorProfileByID(ctx context.Context, authorID string) (*domain.AuthorProfile, error) {
	return u.repo.GetAuthorProfileByID(ctx, authorID, "active")
}

func (u *AuthorProfileUsecases) GetAuthorProfileBySlug(ctx context.Context, slug string) (*domain.AuthorProfile, error) {
	return u.repo.GetAuthorProfileBySlug(ctx, slug, "active")
}

func (u *AuthorProfileUsecases) ListAuthorProfies(ctx context.Context) (*[]domain.AuthorProfile, error) {

	list, err := u.repo.ListAuthorProfies(ctx, "active", "check_not_null")
	if err != nil {
		return nil, err
	}
	return &list, err
}
