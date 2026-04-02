package user

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/delivery/http"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler func(ctx context.Context, evt *messaging.OutboxEvent) error

type EventHandlers struct {
}

type UserModule struct {
	Handler       *http.UserHandler
	EventHandlers *EventHandlers
}

func New(pool *pgxpool.Pool, txManager database.TxManager, outboxRepo outbox.OutboxRepository) *UserModule {

	repo := infrastructure.New(pool)

	authUsecases := application.NewAuthUseCases(txManager, repo, outboxRepo)
	profileUsecases := application.NewProfileUseCases(txManager, repo, outboxRepo)
	notificationUsecases := application.NewNotificationUseCases(txManager, repo, outboxRepo)

	handler := http.New(authUsecases, profileUsecases, notificationUsecases)

	return &UserModule{
		Handler:       handler,
		EventHandlers: nil,
	}
}
