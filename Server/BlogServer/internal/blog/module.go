package blog

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/delivery/event"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/worker"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogWorkerStarter interface {
	StartUpdateRankingTable(ctx context.Context)
	// StartResetBlogDailyViewCount(ctx context.Context)
	// StartResetBlogWeeklyViewCount(ctx context.Context)
}

type Module struct {
	Handler      *http.BlogHandler
	EventHandler *event.EventHandler
	Woker        BlogWorkerStarter
}

// Hide blog module wiring and expose handler, service for other module
func NewBlogModule(pool *pgxpool.Pool, txManager database.TxManager, outboxRepo outbox.OutboxRepository, storageService storage.Storage) *Module {

	repoMapper := infrastructure.NewBlogMapper()
	repo := infrastructure.NewBlogRepository(pool, repoMapper)
	commentRepo := infrastructure.NewCommentRepository(pool, repoMapper)
	blogReactionRepo := infrastructure.NewBlogReactionRepository(pool, repoMapper)
	commentReactionRepo := infrastructure.NewCommentReactionRepository(pool, repoMapper)
	tagRepo := infrastructure.NewTagRepository(pool, repoMapper)

	createBlog := application.NewCreateBlogUseCases(txManager, repo, tagRepo, outboxRepo, storageService)
	getBlog := application.NewGetBlogUseCases(txManager, repo)
	listBlogs := application.NewListBlogsUseCases(txManager, repo)
	deleteBlog := application.NewDeleteBlogUseCases(txManager, repo, outboxRepo)
	commentUsecase := application.NewCommentUseCases(txManager, repo, commentRepo, outboxRepo)
	blogReactionUsecases := application.NewBlogReactionUseCases(txManager, repo, blogReactionRepo, outboxRepo)
	commentReactionUsecases := application.NewCommentReactionUseCases(txManager, commentRepo, commentReactionRepo, outboxRepo)
	blogMetricsUsecases := application.NewBlogMetricsUsecases(txManager, repo)

	handler := http.NewBlogHandler(createBlog, getBlog, listBlogs, deleteBlog, commentUsecase, commentReactionUsecases, blogReactionUsecases, blogMetricsUsecases)

	eventHandler := event.NewEventHandler(txManager, repo, outboxRepo)

	// Worker
	worker := worker.NewBlogWorker(txManager, repo)

	return &Module{
		Handler:      handler,
		EventHandler: eventHandler,
		Woker:        worker,
	}
}
