package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/oklog/ulid/v2"
)

type AuthorIdentityUsecases struct {
	txManager  database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outbox.OutboxRepository
}

func NewAuthorIdentityUsecases(txManager database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outbox.OutboxRepository) *AuthorIdentityUsecases {
	return &AuthorIdentityUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *AuthorIdentityUsecases) CreateAuthor(ctx context.Context, author *domain.AuthorProfile, userID string, createdBy string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		author.AuthorID = ulid.Make().String()
		author.Status = "active"
		err := u.repo.CreateAuthorProfile(ctx, author, userID, createdBy)
		if err != nil {
			return &domain.ErrFailedToCreateAuthorProfile{
				Message: err.Error(),
			}
		}

		event := &domain.AuthorCreatedEvent{
			AuthorID:    author.AuthorID,
			UserID:      userID,
			Slug:        author.Slug,
			DisplayName: author.DisplayName,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *AuthorIdentityUsecases) GetAuthorProfileByID(ctx context.Context, authorID string) (*domain.AuthorProfile, error) {
	return u.repo.GetAuthorProfileByID(ctx, authorID, "active", "check_null")
}

func (u *AuthorIdentityUsecases) GetAuthorProfileBySlug(ctx context.Context, slug string) (*domain.AuthorProfile, error) {
	return u.repo.GetAuthorProfileBySlug(ctx, slug, "active")
}

func (u *AuthorIdentityUsecases) ListAuthorProfiles(ctx context.Context, page int64, limit int64) (*[]domain.AuthorProfile, error) {

	return u.repo.ListAuthorProfiles(ctx, "active", "check_null", page, limit)
}

func (u *AuthorIdentityUsecases) DeleteAuthorProfile(ctx context.Context, authorID string, deletedBy string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		author, err := u.repo.GetAuthorProfileByID(ctx, authorID, "active", "check_null")
		if err != nil {
			return &domain.ErrAuthorNotFound{
				Message: err.Error(),
			}
		}

		if author == nil {
			return &domain.ErrAuthorNotFound{
				Message: "Author not found",
			}
		}

		err = u.repo.DeleteAuthorProfile(ctx, authorID, deletedBy)
		if err != nil {
			return err
		}

		event := &domain.AuthorDeletedEvent{
			AuthorID: authorID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *AuthorIdentityUsecases) HardDeleteAuthorProfile(ctx context.Context, authorID string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		author, err := u.repo.GetAuthorProfileByID(ctx, authorID, "deleted", "check_not_null")
		if err != nil {
			return &domain.ErrAuthorNotFound{
				Message: err.Error(),
			}
		}

		if author == nil {
			return &domain.ErrAuthorNotFound{
				Message: "Author not found",
			}
		}

		err = u.repo.HardDeleteAuthorProfile(ctx, authorID)
		if err != nil {
			return err
		}

		event := &domain.AuthorHardDeletedEvent{
			AuthorID: authorID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *AuthorIdentityUsecases) UpdateAuthorSlug(ctx context.Context, authorID string, slug string, updatedBy string) error {
	return u.repo.UpdateAuthorSlug(ctx, authorID, slug, updatedBy)
}

func (u *AuthorIdentityUsecases) UpdateAuthorStatus(ctx context.Context, authorID string, status string, updatedBy string) error {
	return u.repo.UpdateAuthorStatus(ctx, authorID, status, updatedBy)
}

// Event Handler

func (u *AuthorIdentityUsecases) OnBlogCreated(ctx context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		var evt domain.BlogCountChangedEvent
		err := json.Unmarshal(payload, &evt)
		if err != nil {
			return err
		}

		return u.repo.UpdateAuthorBlogCount(ctx, evt.AuthorID, true)
	})
}
