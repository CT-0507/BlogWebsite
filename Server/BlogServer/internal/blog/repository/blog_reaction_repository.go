package repository

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
)

type BlogReactionRepository interface {
	UpsertBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) (string, string, error)
}
