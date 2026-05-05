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

// func (u *BlogReactionUseCases) CreateBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) error {
// 	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

// 		foundReaction, err := u.blogReactionRepo.GetBlogReaction(ctx, blogReaction.BlogID, blogReaction.UserID)
// 		if err != nil {
// 			return err
// 		}

// 		if foundReaction == nil {
// 			err = u.blogReactionRepo.CreateBlogReaction(ctx, &domain.CreateBlogReaction{
// 				BlogID: blogReaction.BlogID,
// 				UserID: blogReaction.UserID,
// 				Type:   blogReaction.Type,
// 			})
// 			if err != nil {
// 				return err
// 			}

// 			switch blogReaction.Type {
// 			case "like":
// 				err = u.blogRepo.UpdateBlogLikeCount(ctx, blogReaction.BlogID, true)
// 			case "dislike":
// 				err = u.blogRepo.UpdateBlogDislikeCount(ctx, blogReaction.BlogID, true)
// 			default:
// 				return errors.New("Invalid blog reaction type")
// 			}
// 			if err != nil {
// 				return err
// 			}

// 		} else {

// 			if foundReaction.UserID != config.ADMIN_ID && foundReaction.UserID != config.SYSTEM_ID {
// 				return errors.New("Cannot change other user reaction")
// 			}

// 			if foundReaction.Type == blogReaction.Type {
// 				return &contracts.ErrDuplicate{
// 					Message: fmt.Sprintf("Already %s this blog", foundReaction.Type),
// 				}
// 			}

// 			if foundReaction.Type == "like" && blogReaction.Type == "dislike" {
// 				err = u.blogRepo.ChangeLikeToDislike(ctx, blogReaction.BlogID)
// 			} else if foundReaction.Type == "dislike" && blogReaction.Type == "like" {
// 				err = u.blogRepo.ChangeDislikeToLike(ctx, blogReaction.BlogID)
// 			} else {
// 				return errors.New("Action not valid")
// 			}
// 			if err != nil {
// 				return err
// 			}
// 			err = u.blogReactionRepo.UpdateBlogReactionType(ctx, foundReaction.ID, blogReaction.Type)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Save event to outbox table

// 		payload := &map[string]any{}

// 		payloadMarshal, _ := json.Marshal(payload)
// 		err = u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
// 			EventType: "No name",
// 			Payload:   payloadMarshal,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

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
