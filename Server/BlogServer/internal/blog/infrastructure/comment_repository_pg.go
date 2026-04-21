package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		pool: pool,
	}
}

func (r *CommentRepository) CreateComment(c context.Context, newComment *domain.CreateCommentModel) (uuid.UUID, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.CreateComment(c, blogdb.CreateCommentParams{
		BlogID:    newComment.BlogID,
		Content:   newComment.Content,
		ActorType: newComment.ActorType,
		ActorID: pgtype.Text{
			String: utils.GetEmptyStringOnNullStringPtr(newComment.ActorID),
			Valid:  newComment.ActorID != nil,
		},
		Column6:         newComment.RootCommentID == "",
		ParentCommentID: newComment.ParentCommentID,
		Depth:           newComment.Depth,
	})
}

func (r *CommentRepository) GetBlogRootComment(c context.Context, blogID int64) ([]domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.GetBlogRootComment(c, blogID)
	if err != nil {
		return nil, err
	}

	var comments []domain.Comment
	for _, value := range rows {
		v := value
		comments = append(comments, *MapBlogsCommentToComment(&v))
	}
	return comments, nil
}

func (r *CommentRepository) GetChildrenComments(c context.Context, parentCommentID uuid.UUID) ([]domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.GetCommentsByParentComment(c, &parentCommentID)
	if err != nil {
		return nil, err
	}

	var comments []domain.Comment
	for _, value := range rows {
		v := value
		comments = append(comments, *MapBlogsCommentToComment(&v))
	}
	return comments, nil
}

func (r *CommentRepository) GetCommentByID(c context.Context, commentID uuid.UUID) (*domain.Comment, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetCommentByID(c, commentID)
	if err != nil {
		return nil, err
	}

	return MapBlogsCommentToComment(&row), nil
}

func (r *CommentRepository) HideComment(c context.Context, commentID uuid.UUID) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.HideComment(c, commentID)
}

func (r *CommentRepository) DeleteComment(c context.Context, commentID uuid.UUID) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteComment(c, commentID)
}

func (r *CommentRepository) CreateBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.CreateBlogReaction(c, blogdb.CreateBlogReactionParams{
		BlogID: blogReaction.BlogID,
		UserID: blogReaction.UserID,
		Type:   blogReaction.Type,
	})
}

func (r *CommentRepository) SyncBlogReactionCount(c context.Context) error {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.SyncBlogLikeAndDislike(c)
}
