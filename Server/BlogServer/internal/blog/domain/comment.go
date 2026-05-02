package domain

import (
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type Comment struct {
	ID               uuid.UUID  `json:"commentId"`
	BlogID           int64      `json:"blogId"`
	ActorType        string     `json:"actorType"`
	ActorID          *string    `json:"actorId"`
	ActorDisplayName string     `json:"actorDisplayName"`
	ActorAvatarURL   *string    `json:"actorAvatarUrl"`
	Content          string     `json:"content"`
	Status           string     `json:"status"` // active | deleted | hidden
	ParentCommentID  *uuid.UUID `json:"parentCommentId"`
	RootCommentID    uuid.UUID  `json:"rootCommentId"`
	ReplyCount       int64      `json:"replyCount"`
	LikeCount        int32      `json:"likes"`
	DislikeCount     int32      `json:"dislikes"`
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
	ActorID          string
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
