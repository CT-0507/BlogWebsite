package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type FollowerUsecases struct {
	txManager  database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewFollowerUsecases(txManager database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outboxrepo.OutboxRepository) *FollowerUsecases {
	return &FollowerUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *FollowerUsecases) FollowAuthor(ctx context.Context, userID string, authorID string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := u.repo.CreateAuthorFollower(ctx, authorID, userID)
		if err != nil {
			return &domain.ErrFailedToFollowAuthor{
				Message: err.Error(),
			}
		}

		return u.repo.UpdateAuthorFollowerCount(ctx, authorID, true)
	})
}

func (u *FollowerUsecases) UnfollowAuthor(ctx context.Context, userID string, authorID string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := u.repo.DeleteAuthorFollower(ctx, authorID, userID)
		if err != nil {
			return &domain.ErrFailedToUnfollowAuthor{
				Message: err.Error(),
			}
		}

		return u.repo.UpdateAuthorFollowerCount(ctx, authorID, false)
	})
}

func (u *FollowerUsecases) GetAuthorFollowers(ctx context.Context, authorID string, page int64, limit int64) ([]string, error) {
	return u.repo.GetAuthorFollowers(ctx, authorID, page, limit)
}

func (u *FollowerUsecases) GetFollowedAuthors(ctx context.Context, userID string, page int64, limit int64) ([]string, error) {
	return u.repo.GetFollowedAuthors(ctx, userID, page, limit)
}

// func (u *FollowerUsecases) OnAuthorUnfollowed(ctx context.Context, payload []byte) error {
// 	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

// 		var evt domain.AuthorUnfollowedEvent
// 		err := json.Unmarshal(payload, &evt)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}

// 		err = u.repo.DeleteAuthorFollower(ctx, evt.AuthorID, evt.UserID)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}

// 		return u.repo.UpdateAuthorFollowerCount(ctx, evt.AuthorID, false)
// 	})
// }
