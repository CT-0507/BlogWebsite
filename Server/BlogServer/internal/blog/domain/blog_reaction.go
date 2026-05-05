package domain

import (
	"time"

	"github.com/google/uuid"
)

type BlogReaction struct {
	ID        uuid.UUID  `json:"id"`
	BlogID    int64      `json:"blogId"`
	UserID    string     `json:"userId"`
	Type      string     `json:"type"`   // like | dislike
	Status    string     `json:"status"` // active | deleted | hidden
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"-"`
}

type CreateBlogReaction struct {
	BlogID int64  `json:"blogId"`
	UserID string `json:"userId"`
	Type   string `json:"type"` // like | dislike
}
