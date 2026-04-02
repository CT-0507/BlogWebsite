package authors

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler func(ctx context.Context, evt *messaging.OutboxEvent) error

type EventHandlers struct {
	OnAuthorFollowerCountChanged EventHandler
	OnBlogCreated                EventHandler
}

type AuthorsModule struct {
	Handler       *http.AuthorProfileHandler
	EventHandlers *EventHandlers
}

func NewAuthorsModule(pool *pgxpool.Pool, txManager database.TxManager, outboxRepo outbox.OutboxRepository, storageService storage.Storage) *AuthorsModule {
	repo := infrastructure.NewAuthorProfileRepository(pool)

	authorDiscoveryUseCases := application.NewAuthorDiscoveryUsecases(repo)
	authorIdentityUsecases := application.NewAuthorIdentityUsecases(txManager, repo, outboxRepo, storageService)
	authorSocialUsecases := application.NewAuthorSocialUsecases(txManager, repo, outboxRepo)
	authorProfileUsecases := application.NewAuthorProfileUsecases(txManager, repo, outboxRepo)
	authorFollowerUsecases := application.NewFollowerUsecases(txManager, repo, outboxRepo)

	handler := http.NewAuthorProfileHandler(authorDiscoveryUseCases, authorIdentityUsecases, authorSocialUsecases, authorProfileUsecases, authorFollowerUsecases)

	return &AuthorsModule{
		Handler: handler,
		EventHandlers: &EventHandlers{
			OnAuthorFollowerCountChanged: authorFollowerUsecases.OnAuthorFollowerCountChanged,
			OnBlogCreated:                authorIdentityUsecases.OnBlogCreated,
		},
	}
}
