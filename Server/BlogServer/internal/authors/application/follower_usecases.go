package application

import (
	"context"
	"encoding/json"
	"log"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type FollowerUsecases struct {
	txManager  *database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outbox.OutboxRepository
}

func NewFollowerUsecases(txManager *database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outbox.OutboxRepository) *FollowerUsecases {
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
			log.Println(err)
			return &domain.ErrFailedToFollowAuthor{}
		}

		event := &domain.AuthorFollowedEvent{
			AuthorID: authorID,
			UserID:   userID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *FollowerUsecases) UnfollowAuthor(ctx context.Context, userID string, authorID string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		err := u.repo.DeleteAuthorFollower(ctx, authorID, userID)
		if err != nil {
			log.Println(err)
			return &domain.ErrFailedToUnfollowAuthor{}
		}

		event := &domain.AuthorUnfollowedEvent{
			AuthorID: authorID,
			UserID:   userID,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, event.EventName(), payload)
	})
}

func (u *FollowerUsecases) GetAuthorFollowers(ctx context.Context, userID string, authorID string) ([]string, error) {
	return u.repo.GetAuthorFollowers(ctx, authorID, userID)
}

func (u *FollowerUsecases) GetFollowedAuthors(ctx context.Context, userID string) ([]string, error) {
	return u.repo.GetFollowedAuthors(ctx, userID)
}

func (u *FollowerUsecases) OnAuthorFollowed(ctx context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		var evt domain.AuthorFollowedEvent
		err := json.Unmarshal(payload, &evt)
		if err != nil {
			return err
		}

		err = u.repo.UpdateAuthorFollowerCount(ctx, evt.AuthorID)
		if err != nil {
			return err
		}

		newEvt := &domain.UpdateAuthorFollowCountEvent{
			AuthorID: evt.AuthorID,
			UserID:   evt.UserID,
		}

		newEventPayload, err := json.Marshal(newEvt)
		if err != nil {
			return err
		}

		return u.outboxRepo.Insert(ctx, newEvt.EventName(), newEventPayload)
	})
}

func (u *FollowerUsecases) OnAuthorUnfollowed(ctx context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		var evt domain.AuthorUnfollowedEvent
		err := json.Unmarshal(payload, &evt)
		if err != nil {
			return err
		}

		return u.repo.DeleteAuthorFollower(ctx, evt.AuthorID, evt.UserID)

	})
}
