package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommentReaction struct {
	ID        uuid.UUID  `json:"id"`
	CommentID uuid.UUID  `json:"commentId"`
	UserID    string     `json:"userId"`
	Type      string     `json:"type"`   // like | dislike
	Status    string     `json:"status"` // active | deleted | hidden
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

type CreateCommentReaction struct {
	CommentID uuid.UUID `json:"commentId"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"` // like | dislike
}
