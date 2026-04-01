package blog

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/delivery/event"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler      *http.BlogHandler
	EventHandler *event.EventHandler
}

// Hide blog module wiring and expose handler, service for other module
func NewBlogModule(pool *pgxpool.Pool, txManager database.TxManager, outboxRepo outbox.OutboxRepository) *Module {
	repo := infrastructure.NewBlogRepository(pool)

	createBlog := application.NewCreateBlogUseCases(txManager, repo, outboxRepo)
	getBlog := application.NewGetBlogUseCases(txManager, repo)
	listBlogs := application.NewListBlogsUseCases(txManager, repo)
	deleteBlog := application.NewDeleteBlogUseCases(txManager, repo)

	handler := http.NewBlogHandler(createBlog, getBlog, listBlogs, deleteBlog)

	eventHandler := event.NewEventHandler(txManager, repo, outboxRepo)

	return &Module{
		Handler:      handler,
		EventHandler: eventHandler,
	}
}
