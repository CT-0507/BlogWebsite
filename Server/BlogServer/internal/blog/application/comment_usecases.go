package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type CommentUseCases struct {
	txManager   database.TxManager
	blogRepo    repository.BlogRepository
	commentRepo repository.CommentRepository
	outboxRepo  outboxrepo.OutboxRepository
}

func NewCommentUseCases(
	txManager database.TxManager,
	blogRepo repository.BlogRepository,
	commentRepo repository.CommentRepository,
	outboxRepo outboxrepo.OutboxRepository,
) *CommentUseCases {
	return &CommentUseCases{
		txManager:   txManager,
		blogRepo:    blogRepo,
		commentRepo: commentRepo,
		outboxRepo:  outboxRepo,
	}
}

const (
	AUTHOR = "author"
	USER   = "user"
)

func (u *CommentUseCases) CreateComment(c context.Context, newComment *domain.CreateCommentModel, userID string) (*domain.Comment, error) {
	var insertedComment *domain.Comment
	err := u.txManager.WithVoidTx(c, func(ctx context.Context) error {
		var err1 error = nil
		switch newComment.ActorType {
		case USER:
			newComment.ActorID = userID
		case AUTHOR:
			author, err1 := u.blogRepo.GetAuthorProfileByUserID(ctx, userID)
			if err1 != nil {
				return errors.New("Author not found")
			}
			newComment.ActorID = author.AuthorID
			newComment.ActorDisplayName = author.DisplayName
			newComment.ActorAvatarURL = author.AvatarURL
		default:
			return errors.New("Actor type is invalid")
		}

		if newComment.ParentCommentID == nil {
			newComment.RootCommentID = nil
		}

		insertedComment, err1 = u.commentRepo.CreateComment(ctx, &domain.CreateCommentModel{
			BlogID:           newComment.BlogID,
			ActorType:        newComment.ActorType,
			ActorID:          newComment.ActorID,
			ActorDisplayName: newComment.ActorDisplayName,
			ActorAvatarURL:   newComment.ActorAvatarURL,
			Content:          newComment.Content,
			RootCommentID:    newComment.RootCommentID,
			ParentCommentID:  newComment.ParentCommentID,
			Depth:            newComment.Depth,
		})
		if err1 != nil {
			return err1
		}

		// Save event to outbox table

		payload := &map[string]any{}

		payloadMarshal, _ := json.Marshal(payload)
		err1 = u.outboxRepo.Insert(c, &messaging.OutboxEvent{
			EventType: "No name",
			Payload:   payloadMarshal,
		})
		if err1 != nil {
			return err1
		}

		return nil
	})

	return insertedComment, err
}

func (u *CommentUseCases) GetBlogRootComments(c context.Context, blogID int64, userID *string) (int64, []domain.Comment, error) {
	total, err := u.commentRepo.GetBlogRootCommentCount(c, blogID)
	if err != nil {
		return -1, nil, err
	}
	comments, err := u.commentRepo.GetBlogRootComment(c, blogID, userID)
	if err != nil {
		return -1, nil, err
	}
	return total, comments, nil
}

func (u *CommentUseCases) GetChildrenComments(c context.Context, parentCommentID uuid.UUID, userID *string) ([]domain.Comment, error) {
	return u.commentRepo.GetChildrenComments(c, parentCommentID, userID)
}

func (u *CommentUseCases) GetCommentByID(c context.Context, commentID uuid.UUID) (*domain.Comment, error) {
	return u.commentRepo.GetCommentByID(c, commentID)
}

func checkCommentOwnership(actorID *string, userID string) bool {
	if actorID == nil || (*actorID != userID && userID != config.SYSTEM_ID && userID != config.ADMIN_ID) {
		return false
	}
	return true
}

func (u *CommentUseCases) HideComment(c context.Context, commentID uuid.UUID, userID string) (int64, error) {
	comment, err := u.commentRepo.GetCommentByID(c, commentID)
	if err != nil {
		return -1, err
	}
	if comment == nil {
		return -1, errors.New("Comment not Found")
	}

	if checkCommentOwnership(comment.ActorID, userID) {
		return -1, errors.New("User ID not match")
	}
	return u.commentRepo.HideComment(c, commentID)
}

func (u *CommentUseCases) DeleteComment(c context.Context, commentID uuid.UUID, userID string) (int64, error) {
	comment, err := u.commentRepo.GetCommentByID(c, commentID)
	if err != nil {
		return -1, err
	}
	if comment == nil {
		return -1, errors.New("Comment not Found")
	}

	if checkCommentOwnership(comment.ActorID, userID) {
		return -1, errors.New("User ID not match")
	}
	return u.commentRepo.DeleteComment(c, commentID)
}
