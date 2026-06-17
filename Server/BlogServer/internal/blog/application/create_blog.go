package application

import (
	"context"
	"encoding/json"
	"log"

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

		key := storage.GenerateKey("thumbnails", fileParams.FileName)

		uploadResult, err := u.storageService.Save(c, key, fileParams.File, fileParams.ContentType, false)
		if err != nil {
			return nil, err
		}

		// Ensure delete on failure to create user
		defer func() {
			if !success && uploadResult == nil {
				_ = u.storageService.Delete(c, uploadResult.Key)
			}
		}()
		uploadedKey := storage.StripURL(uploadResult.URL)
		blog.ThumbnailUrl = &uploadedKey
	}

	// Process content image src
	var editorData EditorData

	err = json.Unmarshal(blog.ContentJson, &editorData)
	if err != nil {
		return nil, err
	}

	movedURL, err := u.processEditorImages(c, &editorData)
	if err != nil {
		u.rollBackProcessedImages(c, movedURL)
		return nil, err
	}
	defer func() {
		if !success {
			u.rollBackProcessedImages(c, movedURL)
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
			UrlSlug:          insertedBlog.URLSlug,
		}

		newBlog = insertedBlog
		newBlog.ContentJson = updatedJSON
		newBlog.ThumbnailUrl = blog.ThumbnailUrl
		newBlog.Tags = blog.Tags

		newBlog.AuthorID = authorID

		payloadMarshal, _ := json.Marshal(payload)

		// Proceed next step
		err = u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: flows.InceaseAuthorBlogCount,
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

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
	blogID int64,
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

		key := storage.GenerateKey("thumbnails", fileParams.FileName)

		uploadResult, err := u.storageService.Save(c, key,
			fileParams.File,
			fileParams.ContentType,
			false,
		)
		if err != nil {
			return nil, err
		}

		uploadedKey := storage.StripURL(uploadResult.URL)

		newThumbnailURL = &uploadedKey

		// rollback uploaded thumbnail if failed
		defer func() {
			if !success {
				_ = u.storageService.Delete(c, uploadResult.Key)
			}
		}()
	}

	// Parse NEW editor content
	var newEditorData EditorData

	err = json.Unmarshal(payload.ContentJson, &newEditorData)
	if err != nil {
		return nil, err
	}

	// Mark images in editor data to pernament
	movedNewImages, err := u.processEditorImages(c, &newEditorData)
	if err != nil {
		// Rollback successfully marked file
		u.rollBackProcessedImages(c, movedNewImages)
		return nil, err
	}

	// rollback moved images if failed
	defer func() {
		if !success {
			u.rollBackProcessedImages(c, movedNewImages)
		}
	}()

	updatedJSON, err := json.Marshal(newEditorData)
	if err != nil {
		return nil, err
	}

	// Get updated object for handler
	var updatedBlog *domain.Blog

	// For moved image tracking and compensation
	// TODO: Move this logic for background worker
	var proccessedUrls []string
	defer func() {
		if !success && len(proccessedUrls) > 0 {
			u.rollBackProcessedImages(c, proccessedUrls)
		}
	}()
	err = u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		// Edit blog first
		// repo returns before + after
		before, after, err := u.repo.UpdateBlog(ctx, &domain.Blog{
			BlogID:       blogID,
			AuthorID:     payload.AuthorID,
			Title:        payload.Title,
			ContentText:  payload.ContentText,
			ContentJson:  updatedJSON,
			Status:       payload.Status,
			URLSlug:      payload.URLSlug,
			ThumbnailUrl: newThumbnailURL,
		}, newThumbnailURL != nil, userID)
		if err != nil {
			return err
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

		// Thumbnail replacement
		if newThumbnailURL != nil {

			after.ThumbnailUrl = newThumbnailURL

			// old thumbnail cleanup
			if before.ThumbnailUrl != nil &&
				*before.ThumbnailUrl != *newThumbnailURL {
				err = u.storageService.MarkDelete(
					ctx,
					*before.ThumbnailUrl,
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

			err := u.storageService.MarkDelete(ctx, oldURL)
			if err != nil {
				u.rollBackProcessedImages(c, proccessedUrls)
				return err
			}
			proccessedUrls = append(proccessedUrls, oldURL)
		}

		updatedBlog = after
		updatedBlog.Tags = payload.Tags
		updatedBlog.ContentJson = updatedJSON

		return nil
	})

	if err == nil {
		success = true
	} else {
		return nil, err
	}

	return updatedBlog, err
}

func (u *CreateBlogUseCases) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return u.repo.VerifyAuthorIDByUserID(c, userID)
}

func (u *CreateBlogUseCases) UploadTemporaryFile(c context.Context, fileParams storage.FileStorageParams) (string, error) {

	key := storage.GenerateKey("blog_images", fileParams.FileName)

	uploadResult, err := u.storageService.Save(c, key, fileParams.File, fileParams.ContentType, true)
	if err != nil {
		return "", err
	}

	return uploadResult.URL, nil
}

func (u *CreateBlogUseCases) processEditorImages(ctx context.Context, content *EditorData) ([]string, error) {
	var movedURLs []string
	var err error
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

		content.Blocks[i].Data.File.URL = storage.StripURL(url)

		err = u.storageService.MarkPermanent(ctx, url)
		if err != nil {
			break
		}
		movedURLs = append(movedURLs, url)
	}
	return movedURLs, err
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

func (u *CreateBlogUseCases) rollBackProcessedImages(ctx context.Context, movedURLs []string) {
	for _, url := range movedURLs {
		err := u.storageService.MarkPermanent(ctx, url)
		if err != nil {
			log.Println(err)
		}
	}
}
