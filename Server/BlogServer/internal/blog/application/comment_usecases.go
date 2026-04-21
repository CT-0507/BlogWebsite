package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type CommentUseCases struct {
	txManager   database.TxManager
	blogRepo    domain.BlogRepository
	commentRepo domain.CommentRepository
	outboxRepo  outboxrepo.OutboxRepository
}

func NewCommentUseCases(
	txManager database.TxManager,
	blogRepo domain.BlogRepository,
	commentRepo domain.CommentRepository,
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

func (u *CommentUseCases) CreateComment(c context.Context, newComment *domain.CreateCommentModel, userID string) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var actorID string
		switch newComment.ActorType {
		case USER:
			actorID = userID
		case AUTHOR:
			authorID, err := u.blogRepo.VerifyAuthorIDByUserID(ctx, userID)
			if err != nil {
				return errors.New("Author not found")
			}
			actorID = authorID
		default:
			return errors.New("Actor type is invalid")
		}

		if newComment.ParentCommentID != nil {
			newComment.RootCommentID = "-1"
		}

		_, err := u.commentRepo.CreateComment(ctx, &domain.CreateCommentModel{
			BlogID:           newComment.BlogID,
			ActorType:        newComment.ActorType,
			ActorID:          &actorID,
			ActorDisplayName: newComment.ActorDisplayName,
			ActorAvatarURL:   newComment.ActorAvatarURL,
			Content:          newComment.Content,
			RootCommentID:    newComment.RootCommentID,
			ParentCommentID:  newComment.ParentCommentID,
		})
		if err != nil {
			return err
		}

		// Save event to outbox table

		payload := &contracts.BlogCreatedSagaPayload{}

		payloadMarshal, _ := json.Marshal(payload)
		err = u.outboxRepo.Insert(c, &messaging.OutboxEvent{
			EventType: "No name",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *CommentUseCases) GetBlogRootComments(c context.Context, blogID int64) ([]domain.Comment, error) {
	return u.commentRepo.GetBlogRootComment(c, blogID)
}

func (u *CommentUseCases) GetChildrenComments(c context.Context, parentCommentID uuid.UUID) ([]domain.Comment, error) {
	return u.commentRepo.GetChildrenComments(c, parentCommentID)
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

func (u *CommentUseCases) CreateBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		err := u.commentRepo.CreateBlogReaction(ctx, &domain.CreateBlogReaction{
			BlogID: blogReaction.BlogID,
			UserID: blogReaction.UserID,
			Type:   blogReaction.Type,
		})
		if err != nil {
			return err
		}

		// Save event to outbox table

		payload := &contracts.BlogCreatedSagaPayload{}

		payloadMarshal, _ := json.Marshal(payload)
		err = u.outboxRepo.Insert(c, &messaging.OutboxEvent{
			EventType: "No name",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *CommentUseCases) SyncBlogReactionCount(c context.Context) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {
		return u.SyncBlogReactionCount(ctx)
	})
}
