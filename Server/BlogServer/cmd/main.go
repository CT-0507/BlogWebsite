package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/dashboard"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/event_bus"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/notification"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	storage "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/sse"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Application entry point
func main() {

	const PORT = ":8080"

	// Load env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}

	// Uploads
	const RELATIVE_PATH = "../uploads"
	path := os.Getenv("UPLOAD_PATH")
	if dsn == "" {
		path = "." + RELATIVE_PATH
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Create connection pool
	pool, err := database.NewPostgresPool(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Storage
	storage := storage.New(path, "")

	//
	txManager := database.NewTxManager(pool)

	// outbox
	outboxRepo := outbox.New(pool)

	// broker
	broker := sse.NewBroker()

	// sse
	sseHandler := sse.NewSSEHandler(broker)

	// User Module
	userModule := user.New(pool, txManager, outboxRepo)

	// Author Module
	authorModule := authors.NewAuthorsModule(pool, txManager, outboxRepo, storage)

	// Blog CA
	blogModule := blog.NewBlogModule(pool, txManager, outboxRepo)

	// DashBoard
	dashboardHanlder := dashboard.NewDashboardHandler()

	// Notification
	notificationService := notification.NewNotificationService(broker)

	// Register Router
	router := gin.Default()

	// Require authentication
	router.GET("/files/:filepath", func(c *gin.Context) {

		filepath := c.Param("filepath") // /2026/04/03/file.jpg

		fullPath := path + "private" + filepath

		c.File(fullPath)
	})

	// Serve static files
	router.Static(RELATIVE_PATH, path)

	// CORS policy
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "" // allow curl/Postman
		},
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Logger())

	routes.SetupUnprotectedRoutes(router, blogModule.Handler, userModule.Handler, dashboardHanlder, sseHandler, authorModule.Handler)
	routes.SetupProtectedRoutes(router, pool, blogModule.Handler, userModule.Handler, dashboardHanlder, sseHandler, authorModule.Handler)

	saga := saga.NewSagaModule(pool, txManager, outboxRepo)

	bus := event_bus.NewBus()

	// Saga

	// Create blog saga
	bus.Subscribe("create_blog_saga", saga.Orchestrator.StartSaga)
	// Step 1
	bus.Subscribe("CreateBlog", blogModule.EventHandler.CreateBlog)
	bus.Subscribe("InceaseAuthorBlogCount", saga.Orchestrator.HandleEvent)
	bus.Subscribe("CreateBlog.Failed", saga.Orchestrator.HandleFailure)

	// Step 2
	bus.Subscribe("InceaseAuthorBlogCount", authorModule.EventHandler.OnBlogCreated)
	bus.Subscribe("InceaseAuthorBlogCount.Success", saga.Orchestrator.HandleEvent)
	bus.Subscribe("InceaseAuthorBlogCount.Failed", saga.Orchestrator.HandleFailure)
	// Step 2 compensation
	bus.Subscribe("DeleteBlog", blogModule.EventHandler.CreateBlog)
	bus.Subscribe("DeleteBlog.Success", saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe("DeleteBlog.Failed", saga.Orchestrator.HandleCompensationFailure)
	// Notifications
	bus.Subscribe("CreateNotifications", userModule.EventHandler.OnCreateNotifications)

	// Create Author Saga
	bus.Subscribe("create_author_saga", saga.Orchestrator.StartSaga)

	// Step 1
	bus.Subscribe("CreateAuthor", authorModule.EventHandler.OnAuthorFollowerCountChanged)
	bus.Subscribe("CreateAuthor.Failed", saga.Orchestrator.HandleFailure)

	//Step 2
	bus.Subscribe("CreateBlogAuthorCache", blogModule.EventHandler.OnAuthorCreated)
	bus.Subscribe("CreateBlogAuthorCache.Success", saga.Orchestrator.HandleEvent)
	bus.Subscribe("CreateBlogAuthorCache.Failed", saga.Orchestrator.HandleFailure)

	bus.Subscribe("DeleteAuthor", authorModule.EventHandler.OnDeleteAuthor)
	bus.Subscribe("DeleteAuthor.Success", saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe("DeleteAuthor.Failed", saga.Orchestrator.HandleCompensationFailure)

	// bus.Subscribe("blog.created", event_bus.HandlerFunc(authorModule.EventHandlers.OnBlogCreated))
	bus.Subscribe("notification.created", notificationService.PublishNotification)
	bus.Subscribe("authorIdentity.created", blogModule.EventHandler.OnAuthorCreated)
	bus.Subscribe("authorIdentity.deleted", blogModule.EventHandler.OnAuthorDeleted)
	bus.Subscribe("authorIdentity.hardDeleted", blogModule.EventHandler.OnAuthorHardDeleted)
	bus.Subscribe("authorFollower.created", authorModule.EventHandler.OnAuthorFollowerCountChanged)
	bus.Subscribe("authorFollower.deleted", authorModule.EventHandler.OnAuthorFollowerCountChanged)

	worker := outbox.NewOutboxWorker(txManager, bus, outboxRepo)

	go worker.Start(context.Background())

	if err := router.Run(PORT); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
