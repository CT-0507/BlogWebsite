package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthorProfileRepository interface {
	CreateAuthorProfile(c context.Context, author *AuthorProfile, userID uuid.UUID, createdBy uuid.UUID) error
	ListAuthorProfies(c context.Context, status string, deletedCheckMode string) ([]AuthorProfile, error)
}
