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
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/event"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/notification"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/sse"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Application entry point
func main() {

	// Load env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}

	// Create connection pool
	pool, err := database.NewPostgresPool(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// outbox
	outboxRepo := outbox.New(pool)

	// broker
	broker := sse.NewBroker()

	// sse
	sseHandler := sse.NewSSEHandler(broker)

	// User Feature block
	userRepo := user.NewUserRepository()
	userService := user.NewUserService(pool, userRepo)
	userHandler := user.NewUserHandler(userService)

	// Author Module
	authorModule := authors.NewAuthorsModule(pool, outboxRepo)

	// Blog Feature block
	// blogRepo := blog.NewBlogRepository()
	// blogService := blog.NewBlogService(pool, blogRepo, userService, outboxRepo)
	// blogHandler := blog.NewBlogHandler(blogService)

	// Blog CA
	blogModule := blog.NewBlogModule(pool, userService, outboxRepo)

	// DashBoard
	dashboardHanlder := dashboard.NewDashboardHandler()

	// Notification
	notificationService := notification.NewNotificationService(broker)

	// Register Router
	router := gin.Default()

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

	routes.SetupUnprotectedRoutes(router, blogModule.Handler, userHandler, dashboardHanlder, sseHandler, authorModule.Handler)
	routes.SetupProtectedRoutes(router, pool, blogModule.Handler, userHandler, dashboardHanlder, sseHandler, authorModule.Handler)

	bus := event.NewBus()

	bus.Subscribe("blog.created", blogModule.Service.OnBlogPosted)
	bus.Subscribe("blog.created", event.HandlerFunc(authorModule.EventHandlers.OnBlogCreated))
	bus.Subscribe("notification.created", notificationService.PublishNotification)
	bus.Subscribe("authorIdentity.created", blogModule.Service.OnAuthorCreated)
	bus.Subscribe("authorIdentity.deleted", blogModule.Service.OnAuthorDeleted)
	bus.Subscribe("authorIdentity.hardDeleted", blogModule.Service.OnAuthorHardDeleted)
	bus.Subscribe("authorFollower.created", event.HandlerFunc(authorModule.EventHandlers.OnAuthorFollowerCountChanged))
	bus.Subscribe("authorFollower.deleted", event.HandlerFunc(authorModule.EventHandlers.OnAuthorFollowerCountChanged))

	worker := outbox.NewOutboxWorker(pool, bus, outboxRepo)

	go worker.Start(context.Background())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
