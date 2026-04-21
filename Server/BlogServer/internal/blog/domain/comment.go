package domain

import (
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type Comment struct {
	ID               uuid.UUID  `json:"comment_id"`
	BlogID           int64      `json:"blog_id"`
	ActorType        string     `json:"actor_type"`
	ActorID          *string    `json:"actor_id"`
	ActorDisplayName string     `json:"actor_display_name"`
	ActorAvatarURL   *string    `json:"actor_avatar_url"`
	Content          string     `json:"content"`
	Status           string     `json:"status"` // active | deleted | hidden
	ParentCommentID  *uuid.UUID `json:"parent_comment_id"`
	RootCommentID    uuid.UUID  `json:"root_comment_id"`
	Depth            int16      `json:"-"`
	model.Audit
}

type BlogReaction struct {
	ID        uuid.UUID  `json:"blog_reaction_id"`
	BlogID    int64      `json:"blog_id"`
	UserID    string     `json:"user_id"`
	Type      string     `json:"type"`   // like | dislike
	Status    string     `json:"status"` // active | deleted | hidden
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CommentReaction struct {
	ID        uuid.UUID  `json:"id"`
	CommentID uuid.UUID  `json:"comment_id"`
	UserID    string     `json:"user_id"`
	Type      string     `json:"type"`   // like | dislike
	Status    string     `json:"status"` // active | deleted | hidden
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateCommentModel struct {
	BlogID           int64
	ActorType        string
	ActorID          *string
	ActorDisplayName string
	ActorAvatarURL   *string
	Content          string
	ParentCommentID  *uuid.UUID
	RootCommentID    string // Empty on the comment is the root itself
	Depth            int16
}

type CreateBlogReaction struct {
	BlogID int64  `json:"blog_id"`
	UserID string `json:"user_id"`
	Type   string `json:"type"` // like | dislike
}
