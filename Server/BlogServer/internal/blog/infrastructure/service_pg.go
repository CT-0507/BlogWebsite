package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Capsule all blog service
type BlogService struct {
	createBlog           *application.CreateBlogUseCases
	getBlog              *application.GetBlogUseCases
	listBlogs            *application.ListBlogsUseCases
	deleteBlog           *application.DeleteBlogUseCase
	eventHandlerUsecases *application.EventHandlerUsecases
}

func NewBlogService(
	pool *pgxpool.Pool,
	repo domain.BlogRepository,
	userService user.UserService,
	outboxRepo outboxrepo.OutboxRepository,
) application.BlogService {
	txManager := database.NewTxManager(pool)

	return &BlogService{
		createBlog:           application.NewCreateBlogUseCases(txManager, repo, userService, outboxRepo),
		getBlog:              application.NewGetBlogUseCases(txManager, repo),
		listBlogs:            application.NewListBlogsUseCases(txManager, repo),
		deleteBlog:           application.NewDeleteBlogUseCases(txManager, repo),
		eventHandlerUsecases: application.NewEventHandlerUsecases(txManager, repo),
	}
}

func (s *BlogService) CreateBlogStartSaga(c context.Context, blog *domain.Blog, userID string) error {
	return s.createBlog.CreateBlogStartSaga(c, blog, userID)
}

func (s *BlogService) CreateBlog(c context.Context, evt *messaging.OutboxEvent) error {
	return s.createBlog.CreateBlog(c, evt)
}

func (s *BlogService) OnBlogPosted(c context.Context, evt *messaging.OutboxEvent) error {
	return s.createBlog.OnBlogPosted(c, evt)
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

func (s *BlogService) DeleteBlog(ctx context.Context, id int64, userID string) (*int64, error) {
	return s.deleteBlog.DeleteBlog(ctx, id, userID)
}

func (s *BlogService) ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error) {
	return s.listBlogs.ListAuthorBlogsByAuthorID(ctx, authorID)
}

func (s *BlogService) ListAuthorBlogsBySlug(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error) {
	return s.listBlogs.ListAuthorBlogsBySlug(ctx, nickname)
}

func (s *BlogService) OnAuthorCreated(c context.Context, evt *messaging.OutboxEvent) error {
	return s.eventHandlerUsecases.OnAuthorCreated(c, evt)
}

func (s *BlogService) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	return s.createBlog.VerifyAuthorIDByUserID(c, userID)
}

func (s *BlogService) OnAuthorDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return s.eventHandlerUsecases.OnAuthorDeleted(c, evt)
}

func (s *BlogService) OnAuthorHardDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return s.eventHandlerUsecases.OnAuthorHardDeleted(c, evt)
}
