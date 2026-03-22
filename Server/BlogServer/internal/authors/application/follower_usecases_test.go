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

type FollowerUseCasesTestSuite struct {
	suite.Suite
	mockRepo       *mocks.MockAuthorProfileRepository
	txManager      *mocks_test.MockTxManager
	mockOutboxRepo *mocks_test.MockOutboxRepository
	usecases       *application.FollowerUsecases
}

func (s *FollowerUseCasesTestSuite) SetupTest() {
	s.mockRepo = &mocks.MockAuthorProfileRepository{}
	s.txManager = &mocks_test.MockTxManager{}
	s.mockOutboxRepo = &mocks_test.MockOutboxRepository{}
	s.usecases = application.NewFollowerUsecases(
		s.txManager,
		s.mockRepo,
		s.mockOutboxRepo,
	)
}

func TestFollowerUseCasesTestSuite(t *testing.T) {
	suite.Run(t, new(FollowerUseCasesTestSuite))
}

func (s *FollowerUseCasesTestSuite) TestFollowAuthor_Success() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	// Expectation
	s.mockRepo.
		On("CreateAuthorFollower", ctx, authorID, userID).
		Return(nil)

	s.mockOutboxRepo.
		On("Insert", ctx, "authorFollower.created", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(nil)

	// Act
	err := s.usecases.FollowAuthor(ctx, userID, authorID)

	// Assert
	s.NoError(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertExpectations(s.T())
}

func (s *FollowerUseCasesTestSuite) TestFollowAuthorErrorOnInsertingAuthorFollower_Error() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	expectedErr := &domain.ErrFailedToFollowAuthor{
		Message: "",
	}

	// Expectation
	s.mockRepo.
		On("CreateAuthorFollower", ctx, authorID, userID).
		Return(expectedErr)

	// Act
	err := s.usecases.FollowAuthor(ctx, userID, authorID)

	// Assert
	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertNotCalled(s.T(), "Insert")
}

func (s *FollowerUseCasesTestSuite) TestFollowAuthorErrorOnInsertingOutbox_Error() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	expectedErr := errors.New("Error")

	// Expectation
	s.mockRepo.
		On("CreateAuthorFollower", ctx, authorID, userID).
		Return(nil)

	s.mockOutboxRepo.
		On("Insert", ctx, "authorFollower.created", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(expectedErr)

	// Act
	err := s.usecases.FollowAuthor(ctx, userID, authorID)

	// Assert
	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertExpectations(s.T())
}

func (s *FollowerUseCasesTestSuite) TestUnFollowAuthor_Success() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	// Expectation
	s.mockRepo.
		On("DeleteAuthorFollower", ctx, authorID, userID).
		Return(nil)

	s.mockOutboxRepo.
		On("Insert", ctx, "authorFollower.deleted", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(nil)

	// Act
	err := s.usecases.UnfollowAuthor(ctx, userID, authorID)

	// Assert
	s.NoError(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertExpectations(s.T())
}

func (s *FollowerUseCasesTestSuite) TestUnfollowAuthorErrorOnDeletingAuthorFollower_Error() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	expectedErr := &domain.ErrFailedToFollowAuthor{
		Message: "",
	}

	// Expectation
	s.mockRepo.
		On("DeleteAuthorFollower", ctx, authorID, userID).
		Return(expectedErr)

	// Act
	err := s.usecases.UnfollowAuthor(ctx, userID, authorID)

	// Assert
	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertNotCalled(s.T(), "Insert")
}

func (s *FollowerUseCasesTestSuite) TestUnfollowAuthorErrorOnInsertingOutbox_Error() {

	// Arrange

	ctx := context.Background()
	authorID := "author-1"
	userID := "user-1"

	expectedErr := errors.New("Error")

	// Expectation
	s.mockRepo.
		On("DeleteAuthorFollower", ctx, authorID, userID).
		Return(nil)

	s.mockOutboxRepo.
		On("Insert", ctx, "authorFollower.deleted", mock.MatchedBy(func(a []byte) bool {
			return a != nil
		})).
		Return(expectedErr)

	// Act
	err := s.usecases.UnfollowAuthor(ctx, userID, authorID)

	// Assert
	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
	s.mockOutboxRepo.AssertExpectations(s.T())
}
