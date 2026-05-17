package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type BlogReactionUseCases struct {
	txManager        database.TxManager
	blogRepo         repository.BlogRepository
	outboxRepo       outboxrepo.OutboxRepository
	blogReactionRepo repository.BlogReactionRepository
}

func NewBlogReactionUseCases(
	txManager database.TxManager,
	blogRepo repository.BlogRepository,
	blogReactionRepo repository.BlogReactionRepository,
	outboxRepo outboxrepo.OutboxRepository,
) *BlogReactionUseCases {
	return &BlogReactionUseCases{
		txManager:        txManager,
		blogRepo:         blogRepo,
		blogReactionRepo: blogReactionRepo,
		outboxRepo:       outboxRepo,
	}
}

func (u *BlogReactionUseCases) CreateBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) (int, error) {
	var transition repository.ReactionTransition
	err := u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		old, new, err := u.blogReactionRepo.UpsertBlogReaction(ctx, blogReaction)
		if err != nil {
			return err
		}

		switch {
		case old == new:
			return &contracts.ErrDuplicate{
				Message: fmt.Sprintf("Already %s this blog", old)}
		case old == "none" && new == "like":
			transition = repository.AddLike
		case old == "none" && new == "dislike":
			transition = repository.AddDislike
		case old == "like" && new == "dislike":
			transition = repository.LikeToDislike
		case old == "dislike" && new == "like":
			transition = repository.DislikeToLike
		default:
			return fmt.Errorf("invalid reaction transition: %s -> %s", old, new)
		}

		err = u.blogRepo.UpdateBlogReactionCount(ctx, blogReaction.BlogID, transition)
		if err != nil {
			return err
		}

		// Save event to outbox table

		payload := &map[string]any{}

		payloadMarshal, _ := json.Marshal(payload)
		err = u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			EventType: "No name",
			Payload:   payloadMarshal,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return int(transition), err
}

func (u *BlogReactionUseCases) SyncBlogReactionCount(c context.Context) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {
		return u.SyncBlogReactionCount(ctx)
	})
}
