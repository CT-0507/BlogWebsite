package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
)

type AuthorDiscoveryUsecases struct {
	repo domain.AuthorProfileRepository
}

func NewAuthorDiscoveryUsecases(repo domain.AuthorProfileRepository) *AuthorDiscoveryUsecases {
	return &AuthorDiscoveryUsecases{
		repo: repo,
	}
}

func (u *AuthorDiscoveryUsecases) SearchAuthor() {

}

func (u *AuthorDiscoveryUsecases) GetAuthorOwnProfileByUserID(ctx context.Context, userID string) (*domain.AuthorProfile, error) {
	return u.repo.GetAuthorByUserID(ctx, userID)
}
