package repository

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/google/uuid"
)

type CommentRepository interface {
	CreateComment(c context.Context, newComment *domain.CreateCommentModel) (*domain.Comment, error)
	GetBlogRootComment(c context.Context, blogID int64, userID *string) ([]domain.Comment, error)
	GetBlogRootCommentCount(c context.Context, blogID int64) (int64, error)
	GetChildrenComments(c context.Context, parentCommentID uuid.UUID, userID *string) ([]domain.Comment, error)
	GetCommentByID(c context.Context, commentID uuid.UUID) (*domain.Comment, error)
	UpdateComment(c context.Context, commentID uuid.UUID, content, status *string, userID string) (uuid.UUID, error)
	SyncBlogReactionCount(c context.Context) error
	// Metrics
	UpdateCommentReactionCount(c context.Context, commentID uuid.UUID, transition ReactionTransition) error
}
