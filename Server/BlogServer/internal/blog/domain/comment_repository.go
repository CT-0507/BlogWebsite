package domain

import (
	"context"

	"github.com/google/uuid"
)

type CommentRepository interface {
	CreateComment(c context.Context, newComment *CreateCommentModel) (uuid.UUID, error)
	GetBlogRootComment(c context.Context, blogID int64) ([]Comment, error)
	GetChildrenComments(c context.Context, parentCommentID uuid.UUID) ([]Comment, error)
	GetCommentByID(c context.Context, commentID uuid.UUID) (*Comment, error)
	HideComment(c context.Context, commentID uuid.UUID) (int64, error)
	DeleteComment(c context.Context, commentID uuid.UUID) (int64, error)
	CreateBlogReaction(c context.Context, blogReaction *CreateBlogReaction) error
	SyncBlogReactionCount(c context.Context) error
}
