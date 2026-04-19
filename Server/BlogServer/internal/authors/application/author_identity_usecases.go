package application

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
)

type AuthorIdentityUsecases struct {
	txManager      database.TxManager
	repo           domain.AuthorProfileRepository
	outboxRepo     outboxrepo.OutboxRepository
	storageService storage.Storage
}

func NewAuthorIdentityUsecases(
	txManager database.TxManager,
	repo domain.AuthorProfileRepository,
	outboxRepo outboxrepo.OutboxRepository,
	storageService storage.Storage,
) *AuthorIdentityUsecases {
	return &AuthorIdentityUsecases{
		txManager:      txManager,
		repo:           repo,
		outboxRepo:     outboxRepo,
		storageService: storageService,
	}
}

func (u *AuthorIdentityUsecases) CreateAuthor(ctx context.Context, fileParams *domain.CreateUserFileStorageParams, author *domain.AuthorProfile, userID string, createdBy string) error {

	var err error

	if fileParams != nil {
		// Ensure folder on current ymd
		uploadDir, err := utils.EnsureUploadPath("../uploads/temp")
		if err != nil {
			return err
		}

		uploaded := false

		fileParams.FileName = filepath.Join(uploadDir, fileParams.FileName)

		url, err := u.storageService.Upload(fileParams.File, fileParams.FileName, fileParams.ContentType)
		if err != nil {
			return err
		}
		uploaded = true

		// Ensure delete on failure to create user
		defer func() {
			if err != nil && uploaded {
				_ = u.storageService.Delete(fileParams.FileName)
			}
		}()

		author.Avatar = &url
	}

	err = u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		eventPayload := &contracts.AuthorCreatedEventPayload{
			AuthorID:    author.AuthorID,
			UserID:      userID,
			Status:      "active",
			Slug:        author.Slug,
			DisplayName: author.DisplayName,
			SocialLink:  author.SocialLink,
			Email:       author.Email,
			Bio:         author.Bio,
			CreatedBy:   &userID,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.AuthorCreatedEventContext{
			UserID: userID,
			Avatar: author.Avatar,
		}
		context, err := json.Marshal(eventContext)
		if err != nil {
			return err
		}

		sagaID := uuid.New()

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: flows.CreateAuthorSaga,
			SagaID:    &sagaID,
			Payload:   payload,
			Context:   &context,
		})
	})

	return err
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

	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		eventPayload := &contracts.DeleteAuthorKickstartPayload{
			AuthorID:  author.AuthorID,
			UserID:    author.UserID,
			Status:    author.Status,
			Avatar:    author.Avatar,
			UpdatedBy: deletedBy,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		eventContext := &contracts.DeleteAuthorKickstartContext{
			AuthorID: author.AuthorID,
			UserID:   author.UserID,
			Status:   "deleted",
			Avatar:   author.Avatar,
		}

		context, err := json.Marshal(eventContext)
		if err != nil {
			return err
		}

		sagaID := uuid.New()

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    &sagaID,
			EventType: flows.DeleteAuthorSaga,
			Payload:   payload,
			Context:   &context,
		})
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

		event := &contracts.AuthorHardDeletedEvent{
			AuthorID: authorID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType:  event.EventName(),
			Payload:    payload,
			RetryCount: 1,
		})
	})
}

func (u *AuthorIdentityUsecases) UpdateAuthorSlug(ctx context.Context, authorID string, slug string, updatedBy string) error {
	return u.repo.UpdateAuthorSlug(ctx, authorID, slug, updatedBy)
}

func (u *AuthorIdentityUsecases) UpdateAuthorStatus(ctx context.Context, authorID string, status string, updatedBy string) error {
	return u.repo.UpdateAuthorStatus(ctx, authorID, status, updatedBy)
}
