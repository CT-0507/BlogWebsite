package authors

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/delivery/event"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorsModule struct {
	Handler      *http.AuthorProfileHandler
	EventHandler *event.EventHandler
}

func NewAuthorsModule(pool *pgxpool.Pool, txManager database.TxManager, outboxRepo outbox.OutboxRepository, storageService storage.Storage) *AuthorsModule {
	repo := infrastructure.NewAuthorProfileRepository(pool)

	authorDiscoveryUseCases := application.NewAuthorDiscoveryUsecases(repo)
	authorIdentityUsecases := application.NewAuthorIdentityUsecases(txManager, repo, outboxRepo, storageService)
	authorSocialUsecases := application.NewAuthorSocialUsecases(txManager, repo, outboxRepo)
	authorProfileUsecases := application.NewAuthorProfileUsecases(txManager, repo, outboxRepo)
	authorFollowerUsecases := application.NewFollowerUsecases(txManager, repo, outboxRepo)

	handler := http.NewAuthorProfileHandler(authorDiscoveryUseCases, authorIdentityUsecases, authorSocialUsecases, authorProfileUsecases, authorFollowerUsecases)

	eventHandler := event.New(txManager, repo, outboxRepo)

	return &AuthorsModule{
		Handler:      handler,
		EventHandler: eventHandler,
	}
}
