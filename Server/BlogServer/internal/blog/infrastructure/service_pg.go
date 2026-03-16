package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Capsule all blog service
type BlogService struct {
	createBlog *application.CreateBlogUseCases
	getBlog    *application.GetBlogUseCases
	listBlogs  *application.ListBlogsUseCases
	deleteBlog *application.DeleteBlogUseCase
}

func NewBlogService(
	pool *pgxpool.Pool,
	repo domain.BlogRepository,
	userService user.UserService,
	outboxRepo outbox.OutboxRepository,
) application.BlogService {
	txManager := database.NewTxManager(pool)

	return &BlogService{
		createBlog: application.NewCreateBlogUseCases(txManager, repo, userService, outboxRepo),
		getBlog:    application.NewGetBlogUseCases(txManager, repo),
		listBlogs:  application.NewListBlogsUseCases(txManager, repo),
		deleteBlog: application.NewDeleteBlogUseCases(txManager, repo),
	}
}

func (s *BlogService) CreateWithOutBox(c context.Context, blog *domain.Blog) error {
	return s.createBlog.CreateWithOutBox(c, blog)
}

func (s *BlogService) OnBlogPosted(c context.Context, payload []byte) error {
	return s.createBlog.OnBlogPosted(c, payload)
}

func (s *BlogService) GetAll(ctx context.Context) ([]domain.BlogWithAuthorData, error) {
	return s.listBlogs.ListBlogs(ctx)
}

func (s *BlogService) GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error) {
	return s.getBlog.GetBlog(ctx, id)
}

func (s *BlogService) GetBlogByUrlSlug(ctx context.Context, slug string) (*domain.BlogWithAuthorData, error) {
	return s.getBlog.GetBlogByUrlSlug(ctx, slug)
}

func (s *BlogService) DeleteBlog(ctx context.Context, id int64, userID uuid.UUID) (*int64, error) {
	return s.deleteBlog.DeleteBlog(ctx, id, userID)
}

func (s *BlogService) ListAuthorBlogsByAuthorID(ctx context.Context, authorID uuid.UUID) ([]domain.BlogWithAuthorData, error) {
	return s.listBlogs.ListAuthorBlogsByAuthorID(ctx, authorID)
}

func (s *BlogService) ListAuthorBlogsByNickname(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error) {
	return s.listBlogs.ListAuthorBlogsByNickname(ctx, nickname)
}
