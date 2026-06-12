package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	"github.com/google/uuid"
)

type ContactUseCases struct {
	txManager database.TxManager
	repo      domain.UserRepository
}

func NewContactUseCases(txManager database.TxManager, repo domain.UserRepository) *ContactUseCases {
	return &ContactUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

func (u *ContactUseCases) CreateContactForm(c context.Context, contactForm *domain.ContactForm) (*domain.ContactForm, error) {
	var newContact *domain.ContactForm
	err := u.txManager.WithVoidTx(c, func(tx context.Context) error {
		inserted, err := u.repo.CreateContactForm(tx, contactForm)
		if err != nil {
			return err
		}

		newContact = inserted

		notificationContent := map[string]any{
			"content": "You have new contact form",
		}

		contentMarshaled, _ := json.Marshal(notificationContent)

		adminUUID := uuid.MustParse(config.ADMIN_ID)

		// Insert notifcation to admin
		_, err = u.repo.CreateNotification(tx, contentMarshaled, adminUUID, adminUUID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return newContact, nil

}

func (u *ContactUseCases) DeleteContactForm(c context.Context, contactID int64) (int64, error) {

	return u.repo.DeleteContactForm(c, contactID)

}
