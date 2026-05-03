package repository

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
)

type CommentReactionRepository interface {
	UpsertCommentReaction(c context.Context, commentReaction *domain.CreateCommentReaction) (string, string, error)
}
