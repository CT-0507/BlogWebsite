package domain

import (
	"context"

	"github.com/google/uuid"
)

type BlogRepository interface {
	Create(c context.Context, blog *Blog) (*Blog, error)
	FindAll(c context.Context) ([]BlogWithAuthorData, error)
	FindByID(c context.Context, id int64) (*Blog, error)
	// Update(user *Blog) error
	Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error)
}
