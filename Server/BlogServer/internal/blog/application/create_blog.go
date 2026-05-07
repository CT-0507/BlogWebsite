package application

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
)

type CreateBlogUseCases struct {
	txManager      database.TxManager
	repo           repository.BlogRepository
	tagRepo        repository.TagRepository
	storageService storage.Storage
	outboxRepo     outboxrepo.OutboxRepository
}

func NewCreateBlogUseCases(
	txManager database.TxManager,
	repo repository.BlogRepository,
	tagRepo repository.TagRepository,
	outboxRepo outboxrepo.OutboxRepository,
	storageService storage.Storage,
) *CreateBlogUseCases {
	return &CreateBlogUseCases{
		txManager:      txManager,
		repo:           repo,
		tagRepo:        tagRepo,
		outboxRepo:     outboxRepo,
		storageService: storageService,
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
			EventType: flows.CreateBlogSaga,
			Payload:   payloadMarshal,
			Context:   &contextMarshal,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *CreateBlogUseCases) CreateBlog(c context.Context, blog *domain.Blog, userID string, fileParams *storage.FileStorageParams) (*domain.Blog, error) {

	authorID, err := u.repo.VerifyAuthorIDByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	if fileParams != nil {
		// Ensure folder on current ymd
		uploadDir, err := utils.EnsureUploadPath("../uploads")
		if err != nil {
			return nil, err
		}

		uploaded := false
		log.Println("Error here 1")
		fileParams.FileName = filepath.Join(uploadDir, fileParams.FileName)

		url, err := u.storageService.Upload(fileParams.File, fileParams.FileName, fileParams.ContentType)
		if err != nil {
			return nil, err
		}
		uploaded = true

		// Ensure delete on failure to create user
		defer func() {
			if err != nil && uploaded {
				_ = u.storageService.Delete(fileParams.FileName)
			}
		}()
		log.Println("Error here 2")
		blog.ThumbnailUrl = &url
	}

	var newBlog *domain.Blog
	err = u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		blog.AuthorID = authorID

		insertedBlog, err := u.repo.Create(ctx, &domain.Blog{
			AuthorID: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
			Status:   blog.Status,
			URLSlug:  blog.URLSlug,
		})
		if err != nil {
			return err
		}
		log.Println(blog.Tags)
		if len(blog.Tags) > 0 {
			err = u.tagRepo.UpsertTags(ctx, insertedBlog.BlogID, blog.Tags)
			if err != nil {
				return err
			}
		}

		// Save event to outbox table
		payload := &contracts.BlogCreatedEventPayload{
			AuthorID:         authorID,
			BlogID:           blog.BlogID,
			BlogThumbnail:    blog.ThumbnailUrl,
			TruncatedContent: utils.Truncate(blog.Content, 50, true),
			TruncatedTitle:   utils.Truncate(blog.Title, 20, true),
		}

		payloadMarshal, _ := json.Marshal(payload)
		err = u.outboxRepo.Insert(c, &messaging.OutboxEvent{
			EventType: "evt.OnBlogCreated",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		newBlog = insertedBlog

		newBlog.AuthorID = authorID

		return nil
	})

	return newBlog, err
}

func (s *CreateBlogUseCases) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return s.repo.VerifyAuthorIDByUserID(c, userID)
}
