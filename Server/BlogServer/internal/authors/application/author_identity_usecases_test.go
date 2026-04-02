package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application"
	mocks "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application/mocks"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	mocks_test "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthorIdentityUseCasesTestSuite struct {
	suite.Suite
	mockRepo       *mocks.MockAuthorProfileRepository
	txManager      *mocks_test.MockTxManager
	mockOutboxRepo *mocks_test.MockOutboxRepository
	usecases       *application.AuthorIdentityUsecases
}

func (s *AuthorIdentityUseCasesTestSuite) SetupTest() {
	s.mockRepo = &mocks.MockAuthorProfileRepository{}
	s.txManager = &mocks_test.MockTxManager{}
	s.mockOutboxRepo = &mocks_test.MockOutboxRepository{}
	storage := &mocks_test.MockStorage{}
	s.usecases = application.NewAuthorIdentityUsecases(
		s.txManager,
		s.mockRepo,
		s.mockOutboxRepo,
		storage,
	)
}

func TestAuthorIdentityUseCasesTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorIdentityUseCasesTestSuite))
}

func (s *AuthorIdentityUseCasesTestSuite) TestGetAuthorProfileByID_Success() {
	// Arrange

	ctx := context.Background()
	authorID := "123"

	expected := &domain.AuthorProfile{
		AuthorID: "123",
		// Name: "John Doe",
	}

	// Expectation
	s.mockRepo.
		On("GetAuthorProfileByID", ctx, authorID, "active").
		Return(expected, nil)

	// Act
	result, err := s.usecases.GetAuthorProfileByID(ctx, authorID)

	// Assert
	s.NoError(err)
	s.Equal(expected, result)

	s.mockRepo.AssertExpectations(s.T())

}

func (s *AuthorIdentityUseCasesTestSuite) TestGetAuthorProfileByID_Error() {
	// Arrange

	ctx := context.Background()
	authorID := "123"

	expectedErr := errors.New("database error")

	// Expectation
	s.mockRepo.
		On("GetAuthorProfileByID", ctx, authorID, "active").
		Return(nil, expectedErr)

	// Act
	result, err := s.usecases.GetAuthorProfileByID(ctx, authorID)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthorIdentityUseCasesTestSuite) TestGetAuthorProfileBySlug_Success() {
	// Arrange

	ctx := context.Background()
	slug := "123"

	expected := &domain.AuthorProfile{
		AuthorID: "123",
		// Name: "John Doe",
	}

	// Expectation
	s.mockRepo.
		On("GetAuthorProfileBySlug", ctx, slug, "active").
		Return(expected, nil)

	// Act
	result, err := s.usecases.GetAuthorProfileBySlug(ctx, slug)

	// Assert
	s.NoError(err)
	s.Equal(expected, result)

	s.mockRepo.AssertExpectations(s.T())

}

func (s *AuthorIdentityUseCasesTestSuite) TestGetAuthorProfileBySlug_Error() {
	// Arrange

	ctx := context.Background()
	slug := "123"

	expectedErr := errors.New("database error")

	// Expectation
	s.mockRepo.
		On("GetAuthorProfileBySlug", ctx, slug, "active").
		Return(nil, expectedErr)

	// Act
	result, err := s.usecases.GetAuthorProfileBySlug(ctx, slug)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthorIdentityUseCasesTestSuite) TestCreateAuthor_Success() {
	// Arrange

	ctx := context.Background()
	author := &domain.AuthorProfile{}
	userID := "user-1"
	createdBy := userID
	fileParams := &domain.CreateUserFileStorageParams{}
	// Expectation
	s.mockRepo.
		On("CreateAuthorProfile", ctx, mock.MatchedBy(func(a *domain.AuthorProfile) bool {
			return a.AuthorID != ""
		}), userID, createdBy).
		Return(nil)

	s.mockOutboxRepo.
		On("Insert", ctx, "authorIdentity.profileCreated", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(nil)

	// Act
	err := s.usecases.CreateAuthor(ctx, fileParams, author, userID, createdBy)

	// Assert
	s.NoError(err)

	s.mockRepo.AssertExpectations(s.T())
	// s.mockOutboxRepo.AssertExpectations(s.T())
}

func (s *AuthorIdentityUseCasesTestSuite) TestCreateAuthor_ErrorOnInsertingAuthorProfile() {
	// Arrange

	ctx := context.Background()
	author := &domain.AuthorProfile{
		AuthorID: "123",
	}
	userID := "user-1"
	createdBy := userID
	fileParams := &domain.CreateUserFileStorageParams{}

	expectedErr := &domain.ErrFailedToCreateAuthorProfile{
		Message: "Failed to create author profile",
	}

	s.mockRepo.
		On("CreateAuthorProfile", ctx, author, userID, createdBy).
		Return(expectedErr)
	s.mockOutboxRepo.
		On("Insert", ctx, "authorIdentity.profileCreated", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(nil)

	// Act
	err := s.usecases.CreateAuthor(ctx, fileParams, author, userID, createdBy)

	// Assert
	s.Error(err)
	s.Equal(expectedErr, err)
	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertNotCalled(s.T(), "Insert")
}

func (s *AuthorIdentityUseCasesTestSuite) TestCreateAuthor_ErrorOnInsertingOutbox() {
	// Arrange

	ctx := context.Background()
	author := &domain.AuthorProfile{
		AuthorID: "123",
	}
	userID := "user-1"
	createdBy := userID
	fileParams := &domain.CreateUserFileStorageParams{}

	expectedErr := &domain.ErrFailedToCreateAuthorProfile{
		Message: "Failed to create author profile",
	}

	s.mockRepo.
		On("CreateAuthorProfile", ctx, author, userID, createdBy).
		Return(nil)
	s.mockOutboxRepo.
		On("Insert", ctx, "authorIdentity.profileCreated", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(expectedErr)

	// Act
	err := s.usecases.CreateAuthor(ctx, fileParams, author, userID, createdBy)

	// Assert
	s.Error(err)
	s.Equal(expectedErr, err)
	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertExpectations(s.T())
}
