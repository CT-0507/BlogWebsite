package application

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"

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
