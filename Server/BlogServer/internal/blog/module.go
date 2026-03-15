package blog

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler *http.BlogHandler
	Service application.BlogService
}

// Hide blog module wiring and expose handler, service for other module
func NewBlogModule(pool *pgxpool.Pool, userService user.UserService, outboxRepo outbox.OutboxRepository) *Module {
	repo := infrastructure.NewBlogRepository(pool)

	service := infrastructure.NewBlogService(pool, repo, userService, outboxRepo)

	handler := http.NewBlogHandler(service)

	return &Module{
		Handler: handler,
		Service: service,
	}
}
