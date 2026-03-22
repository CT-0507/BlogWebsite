package application_test

import (
	"testing"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application"
	mocks "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application/mocks"
	mocks_test "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/mocks"
	"github.com/stretchr/testify/suite"
)

type AuthorSocialUseCasesTestSuite struct {
	suite.Suite
	mockRepo       *mocks.MockAuthorProfileRepository
	txManager      *mocks_test.MockTxManager
	mockOutboxRepo *mocks_test.MockOutboxRepository
	usecases       *application.AuthorSocialUsecases
}

func (s *AuthorSocialUseCasesTestSuite) SetupTest() {
	s.mockRepo = &mocks.MockAuthorProfileRepository{}
	s.txManager = &mocks_test.MockTxManager{}
	s.mockOutboxRepo = &mocks_test.MockOutboxRepository{}
	s.usecases = application.NewAuthorSocialUsecases(
		s.txManager,
		s.mockRepo,
		s.mockOutboxRepo,
	)
}

func TestAuthorSocialUseCasesTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorIdentityUseCasesTestSuite))
}
