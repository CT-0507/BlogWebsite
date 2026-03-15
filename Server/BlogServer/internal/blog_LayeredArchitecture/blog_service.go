package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog_LayeredArchitecture/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogService interface {
	Create(c context.Context, blog *Blog) error
	CreateWithOutBox(c context.Context, blog *Blog) error
	GetAll(c context.Context) ([]BlogWithAuthorDTO, error)
	GetByID(c context.Context, id int64) (*Blog, error)
	// Update(blog *Blog) error
	Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error)
	OnBlogPosted(c context.Context, payload []byte) error
}

type blogService struct {
	pool        *pgxpool.Pool
	repo        BlogRepository
	userService user.UserService
	outboxRepo  outbox.OutboxRepository
}

func NewBlogService(pool *pgxpool.Pool, repo BlogRepository, userService user.UserService, outboxRepo outbox.OutboxRepository) BlogService {
	return &blogService{
		pool:        pool,
		repo:        repo,
		userService: userService,
		outboxRepo:  outboxRepo,
	}
}

func (s *blogService) withTx(
	ctx context.Context,
	fn func(q *blogdb.Queries) error,
) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	queries := blogdb.New(tx)

	if err := fn(queries); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

//	func (s *blogService) Create(c context.Context, blog *Blog) error {
//		// if blog.Author == "" {
//		// 	return errors.New("name is required")
//		// }
//		ctx, cancel := context.WithTimeout(c, 10*time.Second)
//		defer cancel()
//		q := blogdb.New(s.pool)
//		_, err := s.repo.Create(ctx, q, blog)
//		return err
//	}
func (s *blogService) Create(c context.Context, blog *Blog) error {
	return s.withTx(c, func(q *blogdb.Queries) error {
		_, err := s.repo.Create(c, q, blog)
		return err
	})
}

// Save a box to database and Create an Event to outbox_events table
func (s *blogService) CreateWithOutBox(c context.Context, blog *Blog) error {

	tx, err := s.pool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	queries := blogdb.New(tx)

	blog, err = s.repo.Create(c, queries, blog)
	if err != nil {
		return err
	}

	event := &BlogCreatedEvent{
		BlogID:    blog.BlogID,
		BlogTitle: blog.Title,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = s.outboxRepo.Insert(c, event.EventName(), payload)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (s *blogService) GetAll(c context.Context) ([]BlogWithAuthorDTO, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.FindAll(ctx, q)
}

func (s *blogService) GetByID(c context.Context, id int64) (*Blog, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.FindByID(ctx, q, id)
}

// func (s *blogService) Update(user *Blog) error {
// 	return s.repo.Update(user)
// }

func (s *blogService) Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.Delete(ctx, q, id, userId)
}

func (s *blogService) OnBlogPosted(c context.Context, payload []byte) error {

	tx, err := s.pool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	var evt BlogCreatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	content := fmt.Sprintf("A blog with title %s has just created", evt.BlogTitle)
	not, err := s.userService.CreateNotification(c, content, uuid.MustParse(config.ADMIN_ID), uuid.MustParse(config.SYSTEM_ID))
	if err != nil {
		return err
	}

	// notificationEvent := struct {
	// 	NotifcationID int64
	// 	Content       string
	// 	IsRead        bool
	// }{
	// 	NotifcationID: not.NotificationID,
	// 	Content:       not.Content,
	// 	IsRead:        not.IsRead,
	// }
	notificationPayload, err := json.Marshal(not)
	if err != nil {
		return err
	}

	err = s.outboxRepo.Insert(c, "notification.created", notificationPayload)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}
