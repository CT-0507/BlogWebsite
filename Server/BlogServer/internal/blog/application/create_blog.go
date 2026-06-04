package application

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

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
func (u *CreateBlogUseCases) CreateBlogStartSaga(c context.Context, blog *domain.Blog, userID string) error {

	authorID, err := u.repo.VerifyAuthorIDByUserID(c, userID)
	if err != nil {
		return err
	}
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

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
			// Content:  blog.Content,
			Status: blog.Status,
		}

		payloadMarshal, _ := json.Marshal(payload)
		contextMarshal, _ := json.Marshal(context)
		sagaID := uuid.New()
		err = u.outboxRepo.Insert(c, &messaging.OutboxEvent{
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

	blog.AuthorID = authorID

	success := false
	if fileParams != nil {
		// Ensure folder on current ymd
		uploadDir, err := utils.EnsureUploadPath("../uploads")
		if err != nil {
			return nil, err
		}

		fileParams.FileName = filepath.Join(uploadDir, fileParams.FileName)

		url, err := u.storageService.Upload(fileParams.File, fileParams.FileName, fileParams.ContentType)
		if err != nil {
			return nil, err
		}

		// Ensure delete on failure to create user
		defer func() {
			if !success && url != "" {
				_ = u.storageService.Delete(url)
			}
		}()
		blog.ThumbnailUrl = &url
	}

	// Process content image src
	var editorData EditorData

	err = json.Unmarshal(blog.ContentJson, &editorData)
	if err != nil {
		return nil, err
	}

	movedURL, err := u.processEditorImages(&editorData)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !success {
			u.RollbackMovedImageUrl(movedURL)
		}
	}()

	updatedJSON, err := json.Marshal(editorData)
	if err != nil {
		return nil, err
	}

	// Insert blog
	var newBlog *domain.Blog
	err = u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		insertedBlog, err := u.repo.Create(ctx, &domain.Blog{
			AuthorID:     blog.AuthorID,
			Title:        blog.Title,
			ContentText:  blog.ContentText,
			ContentJson:  updatedJSON,
			Status:       blog.Status,
			URLSlug:      blog.URLSlug,
			ThumbnailUrl: blog.ThumbnailUrl,
		})
		if err != nil {
			return err
		}

		if len(blog.Tags) > 0 {
			err = u.tagRepo.UpsertTags(ctx, insertedBlog.BlogID, blog.Tags)
			if err != nil {
				return err
			}
		}

		// Save event to outbox table
		payload := &contracts.BlogCreatedEventPayload{
			AuthorID:         authorID,
			BlogID:           insertedBlog.BlogID,
			BlogThumbnail:    blog.ThumbnailUrl,
			TruncatedContent: utils.Truncate(blog.ContentText, 50, true),
			TruncatedTitle:   utils.Truncate(blog.Title, 20, true),
		}

		payloadMarshal, _ := json.Marshal(payload)
		err = u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: "evt.OnBlogCreated",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		newBlog = insertedBlog
		newBlog.ContentJson = updatedJSON
		newBlog.ThumbnailUrl = blog.ThumbnailUrl
		newBlog.Tags = blog.Tags

		newBlog.AuthorID = authorID

		return nil
	})

	// Update condition for cleanup job
	if err == nil {
		success = true
	}

	return newBlog, err
}

func (u *CreateBlogUseCases) EditBlog(
	c context.Context,
	blogID string,
	payload *domain.Blog,
	userID string,
	fileParams *storage.FileStorageParams,
) (*domain.Blog, error) {

	authorID, err := u.repo.VerifyAuthorIDByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	payload.AuthorID = authorID

	success := false

	// Thumbnail upload handling
	var newThumbnailURL *string

	if fileParams != nil {

		uploadDir, err := utils.EnsureUploadPath("../uploads")
		if err != nil {
			return nil, err
		}

		fileParams.FileName = filepath.Join(uploadDir, fileParams.FileName)

		url, err := u.storageService.Upload(
			fileParams.File,
			fileParams.FileName,
			fileParams.ContentType,
		)
		if err != nil {
			return nil, err
		}

		newThumbnailURL = &url

		// rollback uploaded thumbnail if failed
		defer func() {
			if !success {
				_ = u.storageService.Delete(url)
			}
		}()
	}

	// Parse NEW editor content
	var newEditorData EditorData

	err = json.Unmarshal(payload.ContentJson, &newEditorData)
	if err != nil {
		return nil, err
	}

	// Move newly uploaded temp images
	movedNewImages, err := u.processEditorImages(&newEditorData)
	if err != nil {
		return nil, err
	}

	// rollback moved images if failed
	defer func() {
		if !success {
			u.RollbackMovedImageUrl(movedNewImages)
		}
	}()

	updatedJSON, err := json.Marshal(newEditorData)
	if err != nil {
		return nil, err
	}

	var updatedBlog *domain.Blog
	err = u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		// Edit blog first
		// repo returns before + after
		before, after, err := u.repo.UpdateBlog(ctx, &domain.Blog{
			AuthorID:     payload.AuthorID,
			Title:        payload.Title,
			ContentText:  payload.ContentText,
			ContentJson:  updatedJSON,
			Status:       payload.Status,
			URLSlug:      payload.URLSlug,
			ThumbnailUrl: newThumbnailURL,
		}, userID)
		if err != nil {
			return err
		}
		// Thumbnail replacement
		if newThumbnailURL != nil {

			after.ThumbnailUrl = newThumbnailURL

			// old thumbnail cleanup
			if before.ThumbnailUrl != nil &&
				*before.ThumbnailUrl != *newThumbnailURL {

				tempThumb := utils.SwapTemp(*before.ThumbnailUrl, true)

				err = u.storageService.MoveFile(
					*before.ThumbnailUrl,
					tempThumb,
				)
				if err != nil {
					return err
				}
			}
		}

		// Compare old/new editor images
		var beforeEditor EditorData
		if err := json.Unmarshal(before.ContentJson, &beforeEditor); err != nil {
			return err
		}
		oldImages := u.extractEditorImageURLs(&beforeEditor)
		newImages := u.extractEditorImageURLs(&newEditorData)

		removedImages := u.difference(oldImages, newImages)

		// move removed images into temp folder
		for _, oldURL := range removedImages {

			// skip temp images
			if strings.Contains(oldURL, "/uploads/temp/") {
				continue
			}

			tempURL := utils.SwapTemp(oldURL, true)

			err := u.storageService.MoveFile(oldURL, tempURL)
			if err != nil {
				return err
			}
		}

		// tags
		if len(payload.Tags) > 0 {
			err = u.tagRepo.UpsertTags(
				ctx,
				after.BlogID,
				payload.Tags,
			)
			if err != nil {
				return err
			}
		}
		updatedBlog = after
		updatedBlog.ContentJson = updatedJSON
		updatedBlog.Tags = payload.Tags

		return nil
	})

	if err == nil {
		success = true
	}
	if err != nil {
		return nil, err
	}

	return updatedBlog, err
}

func (u *CreateBlogUseCases) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return u.repo.VerifyAuthorIDByUserID(c, userID)
}

func (u *CreateBlogUseCases) SaveBlogImageToTempFolder(c context.Context, fileParams storage.FileStorageParams) (string, error) {

	// Ensure folder on current ymd
	uploadDir, err := utils.EnsureUploadPath("../uploads/temp")
	if err != nil {
		return "", err
	}

	fileParams.FileName = filepath.Join(uploadDir, fileParams.FileName)

	url, err := u.storageService.Upload(fileParams.File, fileParams.FileName, fileParams.ContentType)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *CreateBlogUseCases) processEditorImages(content *EditorData) ([]string, error) {
	var movedURLs []string
	for i, block := range content.Blocks {
		if block.Type != "image" {
			continue
		}

		// safety check
		if block.Data.File == nil {
			continue
		}

		url := block.Data.File.URL

		if url == "" {
			continue
		}

		if strings.Contains(url, "/uploads/temp/") {

			dst := utils.SwapTemp(url, false)

			err := u.storageService.MoveFile(url, dst)
			if err != nil {
				return nil, err
			}

			content.Blocks[i].Data.File.URL = dst

			movedURLs = append(movedURLs, dst)
		}
	}
	return movedURLs, nil
}

func (u *CreateBlogUseCases) extractEditorImageURLs(content *EditorData) []string {

	var urls []string

	for _, block := range content.Blocks {

		if block.Type != "image" {
			continue
		}

		if block.Data.File == nil {
			continue
		}

		url := block.Data.File.URL

		if url == "" {
			continue
		}

		urls = append(urls, url)
	}

	return urls
}

func (u *CreateBlogUseCases) difference(oldImages, newImages []string) []string {

	newMap := make(map[string]bool)

	for _, v := range newImages {
		newMap[v] = true
	}

	var removed []string

	for _, v := range oldImages {
		if !newMap[v] {
			removed = append(removed, v)
		}
	}

	return removed
}

func (u *CreateBlogUseCases) RollbackMovedImageUrl(movedURLs []string) {
	for _, url := range movedURLs {

		if !strings.Contains(url, "/uploads/") {
			continue
		}

		dst := utils.SwapTemp(url, true)

		_ = u.storageService.MoveFile(url, dst)
	}
}
