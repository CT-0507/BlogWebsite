package infrastructure

import (
	"context"
	"errors"
	"slices"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewCommentRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *CommentRepository {
	return &CommentRepository{
		pool:   pool,
		mapper: mapper,
	}
}

func (r *CommentRepository) CreateComment(c context.Context, newComment *domain.CreateCommentModel) (*domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	var rootCommentID *uuid.UUID = &uuid.Nil
	if newComment.RootCommentID != nil {
		rootCommentID = newComment.RootCommentID
	}

	inserted, err := q.CreateComment(c, blogdb.CreateCommentParams{
		BlogID:    newComment.BlogID,
		Content:   newComment.Content,
		ActorType: newComment.ActorType,
		ActorID: pgtype.Text{
			String: newComment.ActorID,
			Valid:  true,
		},
		ActorDisplayName: newComment.ActorDisplayName,
		RootCommentID:    *rootCommentID,
		ParentCommentID:  newComment.ParentCommentID,
		Depth:            newComment.Depth,
	})
	if utils.IsDuplicateKey(err) {
		return nil, &contracts.ErrDuplicate{
			Message: err.Error(),
		}
	}
	return r.mapper.MapBlogsCommentToComment(&inserted), err
}

func (r *CommentRepository) GetBlogRootComment(c context.Context, blogID int64, userID *string) ([]domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)
	var comments []domain.Comment
	if userID != nil {
		rows, err := q.GetBlogRootCommentWithUserReaction(c, blogdb.GetBlogRootCommentWithUserReactionParams{
			BlogID: blogID,
			UserID: *userID,
		})
		if err != nil {
			return nil, err
		}
		for _, value := range rows {
			v := value
			mappedV, err := r.mapper.MapBlogRootCommentWithReactionRow(&v)
			if err != nil {
				return nil, err
			}
			comments = append(comments, *mappedV)
		}
	} else {
		rows, err := q.GetBlogRootComment(c, blogID)
		if err != nil {
			return nil, err
		}
		for _, value := range rows {
			v := value
			mappedV, err := r.mapper.MapBlogRootCommentRow(&v)
			if err != nil {
				return nil, err
			}
			comments = append(comments, *mappedV)
		}
	}

	return comments, nil
}

func (r *CommentRepository) GetBlogRootCommentCount(c context.Context, blogID int64) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.GetBlogRootCommentCount(c, blogID)
}

func (r *CommentRepository) GetChildrenComments(c context.Context, parentCommentID uuid.UUID, userID *string) ([]domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	if userID != nil {
		rows, err := q.GetCommentsByParentCommentUserWithReaction(c, blogdb.GetCommentsByParentCommentUserWithReactionParams{
			ParentCommentID: &parentCommentID,
			UserID:          *userID,
		})
		if err != nil {
			return nil, err
		}
		var comments []domain.Comment
		for _, value := range rows {
			v := value
			mappedV, err := r.mapper.MapCommentsByParentCommentWithReactionRow(&v)
			if err != nil {
				return nil, err
			}
			comments = append(comments, *mappedV)
		}
		return comments, nil
	} else {
		rows, err := q.GetCommentsByParentComment(c, &parentCommentID)
		if err != nil {
			return nil, err
		}

		var comments []domain.Comment
		for _, value := range rows {
			v := value
			mappedV, err := r.mapper.MapCommentsByParentCommentRow(&v)
			if err != nil {
				return nil, err
			}
			comments = append(comments, *mappedV)
		}
		return comments, nil
	}

}

func (r *CommentRepository) GetCommentByID(c context.Context, commentID uuid.UUID) (*domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetCommentByID(c, commentID)
	if err != nil {
		return nil, err
	}

	return r.mapper.MapBlogsCommentToComment(&row), nil
}

func (r *CommentRepository) UpdateComment(c context.Context, commentID uuid.UUID, content, status *string, userID string) (uuid.UUID, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	if status != nil {
		validStatuses := []string{"deleted", "active", "hidden"}
		if !slices.Contains(validStatuses, *status) {
			return uuid.Nil, errors.New("Invalid status")
		}
	}

	return q.UpdateComment(c, blogdb.UpdateCommentParams{
		Content: utils.GetTextTypeFromNullableString(content),
		Status:  utils.GetTextTypeFromNullableString(status),
		ID:      commentID,
		ActorID: pgtype.Text{
			String: userID,
			Valid:  true,
		},
		IsAdmin: false,
	})
}

func (r *CommentRepository) SyncBlogReactionCount(c context.Context) error {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.SyncBlogLikeAndDislike(c)
}

func (r *CommentRepository) UpdateCommentReactionCount(c context.Context, commentID uuid.UUID, transition repository.ReactionTransition) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	var likeDelta int32 = 0
	var dislikeDelta int32 = 0

	switch transition {
	case repository.AddLike:
		likeDelta++
	case repository.AddDislike:
		dislikeDelta++
	case repository.LikeToDislike:
		likeDelta--
		dislikeDelta++
	case repository.DislikeToLike:
		likeDelta++
		dislikeDelta--
	}

	return q.UpdateCommentReactionCount(c, blogdb.UpdateCommentReactionCountParams{
		LikeCount:    likeDelta,
		DislikeCount: dislikeDelta,
		ID:           commentID,
	})
}
