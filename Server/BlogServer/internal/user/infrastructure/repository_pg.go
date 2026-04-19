package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/infrastructure/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(c context.Context, user *domain.User) (*domain.User, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	newUser, err := q.CreateUser(c, userdb.CreateUserParams{
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Roles[0],
	})
	if err != nil {
		return nil, err
	}

	return UserDTOToUser(&newUser), nil
}

func (r *UserRepository) CountByUsername(c context.Context, username string) (int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	return q.CountUserWithEmail(c, username)
}

func (r *UserRepository) GetUserByUsername(c context.Context, username string) (*domain.User, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	user, err := q.GetUserByUsername(c, username)
	if err != nil {
		return nil, err
	}
	return UserDTOToUser(&user), nil
}

func (r *UserRepository) UpdateLastLogout(c context.Context, userID uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	return q.UpdateLastLogout(c, userID)
}

func (r *UserRepository) GetUserByID(c context.Context, userID uuid.UUID) (*domain.User, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	user, err := q.GetUserByID(c, userID)
	if err != nil {
		return nil, err
	}
	return UserDTOToUser(&user), nil
}

func (r *UserRepository) UpdateEmail(c context.Context, userID uuid.UUID, email string, updatedBy *uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	_, err := q.UpdateUserEmail(c, userdb.UpdateUserEmailParams{
		UserID: userID,
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
		UpdatedBy: updatedBy,
	})
	return err
}

func (r *UserRepository) UpdateData(c context.Context, userID uuid.UUID, user *domain.User, updatedBy *uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	_, err := q.UpdateUserData(c, userdb.UpdateUserDataParams{
		UserID:    userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UpdatedBy: updatedBy,
	})
	return err
}

func (r *UserRepository) UpdatePassword(c context.Context, userID uuid.UUID, hashedPassword string, updatedBy *uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	_, err := q.UpdateUserPassword(c, userdb.UpdateUserPasswordParams{
		UserID:    userID,
		Password:  hashedPassword,
		UpdatedBy: updatedBy,
	})
	return err
}

func (r *UserRepository) GetNotificationsByUserID(c context.Context, userID uuid.UUID) ([]domain.Notification, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	rows, err := q.GetUserNotiticationsByID(c)
	if err != nil {
		return nil, err
	}
	var notifications []domain.Notification
	for _, value := range rows {
		v := value
		notifications = append(notifications, *NotificationDTOToNotification(&v))
	}
	return notifications, nil
}

func (r *UserRepository) CreateNotification(c context.Context, content []byte, userID uuid.UUID, createdBy uuid.UUID) (*domain.Notification, error) {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	notification, err := q.CreateNotification(c, userdb.CreateNotificationParams{
		UserID:    &userID,
		Content:   content,
		CreatedBy: &createdBy,
	})
	return NotificationDTOToNotification(&notification), err
}

func (r *UserRepository) CreateNotifications(c context.Context, nots []domain.Notification) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	systemUUID := uuid.MustParse(config.SYSTEM_ID)

	var params []userdb.CreateNotificationsParams
	for _, value := range nots {
		v := value
		uuid := uuid.MustParse(value.UserID)
		params = append(params, userdb.CreateNotificationsParams{
			UserID:    &uuid,
			Content:   v.Content,
			CreatedBy: &systemUUID,
		})
	}

	_, err := q.CreateNotifications(c, params)
	return err
}
func (r *UserRepository) UpdateNotificationByID(c context.Context, notificationID int64, status bool, updatedBy *uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	return q.UpdateNotification(c, userdb.UpdateNotificationParams{
		IsRead:         status,
		NotificationID: notificationID,
		UpdatedBy:      updatedBy,
	})
}

func (r *UserRepository) MarkUserAsDeleted(c context.Context, userID uuid.UUID, updatedBy uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	return q.MarkUserAsDeleted(c, userdb.MarkUserAsDeletedParams{
		UserID:    userID,
		UpdatedBy: &updatedBy,
	})
}

func (r *UserRepository) RestoreUserByID(c context.Context, userID uuid.UUID, updatedBy uuid.UUID) error {

	db := utils.GetExecutor(c, r.pool)

	q := userdb.New(db)

	return q.RestoreUserByID(c, userdb.RestoreUserByIDParams{
		UserID:    userID,
		UpdatedBy: &updatedBy,
	})
}
