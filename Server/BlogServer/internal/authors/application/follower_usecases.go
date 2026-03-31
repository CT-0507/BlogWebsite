package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
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

		event := &contracts.AuthorFollowedEvent{
			AuthorID:    authorID,
			UserID:      userID,
			IsIncrement: true,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType:  event.EventName(),
			Payload:    payload,
			RetryCount: 1,
		})
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

		event := &contracts.AuthorUnfollowedEvent{
			AuthorID:    authorID,
			UserID:      userID,
			IsIncrement: false,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType:  event.EventName(),
			Payload:    payload,
			RetryCount: 1,
		})
	})
}

func (u *FollowerUsecases) GetAuthorFollowers(ctx context.Context, authorID string, page int64, limit int64) ([]string, error) {
	return u.repo.GetAuthorFollowers(ctx, authorID, page, limit)
}

func (u *FollowerUsecases) GetFollowedAuthors(ctx context.Context, userID string, page int64, limit int64) ([]string, error) {
	return u.repo.GetFollowedAuthors(ctx, userID, page, limit)
}

func (u *FollowerUsecases) OnAuthorFollowerCountChanged(ctx context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		var event contracts.FollowCountChangedEvent
		err := json.Unmarshal(evt.Payload, &event)
		if err != nil {
			return err
		}

		err = u.repo.UpdateAuthorFollowerCount(ctx, event.AuthorID, event.IsIncrement)
		if err != nil {
			return err
		}

		event2 := &contracts.FollowCountChangedEvent{
			AuthorID:    event.AuthorID,
			UserID:      event.UserID,
			IsIncrement: event.IsIncrement,
		}

		eventPayload, _ := json.Marshal(event2)

		return u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType:  event2.EventName(),
			Payload:    eventPayload,
			RetryCount: 1,
		})
	})
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
