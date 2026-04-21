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
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
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
	bus.Subscribe(flows.CreateBlogSaga, saga.Orchestrator.StartSaga)
	// Step 1
	bus.Subscribe(flows.CreateBlog, blogModule.EventHandler.CreateBlog)
	bus.Subscribe(flows.CreateBlogSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.CreateBlogFailed, saga.Orchestrator.HandleFailure)

	// Step 2
	bus.Subscribe(flows.InceaseAuthorBlogCount, authorModule.EventHandler.OnBlogCreated)
	bus.Subscribe(flows.InceaseAuthorBlogCountSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.InceaseAuthorBlogCountFailed, saga.Orchestrator.HandleFailure)
	// Step 1 compensation
	bus.Subscribe(flows.CreateBlogCompensation, blogModule.EventHandler.OnCreateBlogCompensation)
	bus.Subscribe(flows.CreateBlogCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.CreateBlogCompensationFailed, saga.Orchestrator.HandleCompensationFailure)
	// Notifications
	bus.Subscribe("CreateNotifications", userModule.EventHandler.OnCreateNotifications)

	// Create Author Saga
	bus.Subscribe(flows.CreateAuthorSaga, saga.Orchestrator.StartSaga)

	// Step 1
	bus.Subscribe(flows.CreateAuthor, authorModule.EventHandler.OnAuthorCreate)
	bus.Subscribe(flows.CreateAuthorSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.CreateAuthorFailed, saga.Orchestrator.HandleFailure)

	//Step 2
	bus.Subscribe(flows.CreateBlogAuthorCache, blogModule.EventHandler.OnAuthorCreated)
	bus.Subscribe(flows.CreateBlogAuthorCacheSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.CreateBlogAuthorCacheFailed, saga.Orchestrator.HandleFailure)

	bus.Subscribe(flows.CreateAuthorCompensation, authorModule.EventHandler.OnCreateAuthorCompensation)
	bus.Subscribe(flows.CreateAuthorCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.CreateAuthorCompensationFailed, saga.Orchestrator.HandleCompensationFailure)

	// Delete Author Saga
	bus.Subscribe(flows.DeleteAuthorSaga, saga.Orchestrator.StartSaga)

	// Step 1: Mark author as deleted
	bus.Subscribe(flows.DeleteAuthor, authorModule.EventHandler.OnDeleteAuthor)
	bus.Subscribe(flows.DeleteAuthorSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.DeleteAuthorFailed, saga.Orchestrator.HandleFailure)

	// Step 2
	bus.Subscribe(flows.DeleteBlogAuthorCache, blogModule.EventHandler.OnDeleteBlogAuthorCache)
	bus.Subscribe(flows.DeleteBlogAuthorCacheSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.DeleteBlogAuthorCacheFailed, saga.Orchestrator.HandleFailure)

	bus.Subscribe(flows.DeleteAuthorCompensation, authorModule.EventHandler.OnDeleteAuthorCompensation)
	bus.Subscribe(flows.DeleteAuthorCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.DeleteAuthorCompensationFailed, saga.Orchestrator.HandleCompensationFailure)

	// Delete user saga
	bus.Subscribe(flows.DeleteUserSaga, saga.Orchestrator.StartSaga)

	// Step 1
	bus.Subscribe(flows.DeleteUser, userModule.EventHandler.OnDeleteUser)
	bus.Subscribe(flows.DeleteUserSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.DeleteUserFailed, saga.Orchestrator.HandleFailure)

	// Step 2
	bus.Subscribe(flows.CleanUpAuthorProfile, authorModule.EventHandler.OnUserDeleted)
	bus.Subscribe(flows.CleanUpAuthorProfileSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.CleanUpAuthorProfileFailed, saga.Orchestrator.HandleFailure)
	// Step 1 Compensation
	bus.Subscribe(flows.DeleteUserCompensation, authorModule.EventHandler.OnUserDeletedCompensation)
	bus.Subscribe(flows.DeleteUserCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.DeleteUserCompensationFailed, saga.Orchestrator.HandleCompensationFailure)

	// Step Compensation
	bus.Subscribe(flows.CleanUpAuthorProfileCompensation, authorModule.EventHandler.OnUserDeletedCompensation)
	bus.Subscribe(flows.CleanUpAuthorProfileCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.CleanUpAuthorProfileCompensationFailed, saga.Orchestrator.HandleCompensationFailure)

	// Delete Blog Saga
	bus.Subscribe(flows.DeleteBlogSaga, saga.Orchestrator.StartSaga)

	// Step 1
	bus.Subscribe(flows.DeleteBlog, blogModule.EventHandler.OnDeleteBlog)
	bus.Subscribe(flows.DeleteBlogSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.DeleteBlogFailed, saga.Orchestrator.HandleFailure)

	// Step 2
	bus.Subscribe(flows.DecreaseAuthorBlogCount, authorModule.EventHandler.OnDecreaseAuthorBlogCount)
	bus.Subscribe(flows.DecreaseAuthorBlogCountSuccess, saga.Orchestrator.HandleEvent)
	bus.Subscribe(flows.DecreaseAuthorBlogCountFailed, saga.Orchestrator.HandleFailure)

	// Step 1 compensation
	bus.Subscribe(flows.DeleteBlogCompensation, blogModule.EventHandler.OnDeleteBlogCompensation)
	bus.Subscribe(flows.DeleteBlogCompensationSuccess, saga.Orchestrator.HandleCompensationSuccess)
	bus.Subscribe(flows.DeleteBlogCompensationFailed, saga.Orchestrator.HandleCompensationFailure)

	// bus.Subscribe("blog.created", event_bus.HandlerFunc(authorModule.EventHandlers.OnBlogCreated))
	bus.Subscribe("notification.created", notificationService.PublishNotification)

	worker := outbox.NewOutboxWorker(txManager, bus, outboxRepo)

	go worker.Start(context.Background())

	if err := router.Run(PORT); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
