package application

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
)

type GetBlogUseCases struct {
	txManager database.TxManager
	repo      repository.BlogRepository
}

func NewGetBlogUseCases(txManager database.TxManager, repo repository.BlogRepository) *GetBlogUseCases {
	return &GetBlogUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

func (u *GetBlogUseCases) GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error) {
	blog, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if blog.ThumbnailUrl != nil {
		thumbnailWithDomain, err := storage.AddDomain(*blog.ThumbnailUrl)
		if err != nil {
			return nil, err
		}
		blog.ThumbnailUrl = &thumbnailWithDomain
	}

	contentJson, err := u.proccessUrlToDomain(blog.ContentJson)
	if err != nil {
		return nil, err
	}
	blog.ContentJson = contentJson

	return blog, nil
}

func (u *GetBlogUseCases) proccessUrlToDomain(contentJson json.RawMessage) (json.RawMessage, error) {

	var editorData EditorData

	err := json.Unmarshal(contentJson, &editorData)
	if err != nil {
		return nil, err
	}

	for i, block := range editorData.Blocks {
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

		urlWithDomain, err := storage.AddDomain(url)
		if err != nil {
			return nil, err
		}

		editorData.Blocks[i].Data.File.URL = urlWithDomain

	}

	rawMessage, err := json.Marshal(editorData)
	if err != nil {
		return nil, err
	}

	return rawMessage, nil
}

func (u *GetBlogUseCases) GetBlogByUrlSlug(ctx context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error) {
	result, err := u.repo.FindByUrlSlug(ctx, slug, userID)

	if err != nil {
		return nil, err
	}

	if len(result.Tags) == 0 {
		result.Tags = []string{}
	}

	if result.ThumbnailUrl != nil {
		thumbnailWithDomain, err := storage.AddDomain(*result.ThumbnailUrl)
		if err != nil {
			return nil, err
		}
		result.ThumbnailUrl = &thumbnailWithDomain
	}

	contentJson, err := u.proccessUrlToDomain(result.ContentJson)
	if err != nil {
		return nil, err
	}
	result.ContentJson = contentJson

	go func(id int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := u.repo.UpdateViewCount(ctx, id); err != nil {
			log.Println("daily metric failed:", err)
		}
	}(result.BlogID)

	return result, err
}
