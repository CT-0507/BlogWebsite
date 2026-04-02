package application

import (
	"context"

	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	"github.com/google/uuid"
)

type NotificationUseCases struct {
	txManager  database.TxManager
	repo       domain.UserRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewNotificationUseCases(txManager database.TxManager, repo domain.UserRepository, outboxRepo outboxrepo.OutboxRepository) *NotificationUseCases {
	return &NotificationUseCases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *NotificationUseCases) GetUserNotifications(c context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	return u.repo.GetNotificationsByUserID(c, userID)
}

func (u *NotificationUseCases) UpdateNotificationStatus(c context.Context, notID int64, status bool, updatedBy *uuid.UUID) error {
	return u.repo.UpdateNotificationByID(c, notID, status, updatedBy)
}

func (u *NotificationUseCases) CreateNotification(c context.Context, content string, userID uuid.UUID, createdBy uuid.UUID) (*domain.Notification, error) {
	systemId := uuid.MustParse(config.SYSTEM_ID)
	not, err := u.repo.CreateNotification(c, content, userID, systemId)
	if err != nil {
		return nil, err
	}
	return not, nil
}
