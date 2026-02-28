package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/dashboard"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to reach database server: %v", err)
	}

	// Blog Feature block
	blogRepo := blog.NewBlogRepository()
	blogService := blog.NewBlogService(pool, blogRepo)
	blogHandler := blog.NewBlogHandler(blogService)

	// User Feature block
	userRepo := user.NewUserRepository()
	userService := user.NewUserService(pool, userRepo)
	userHandler := user.NewUserHandler(userService)

	// DashBoard
	dashboardHanlder := dashboard.NewDashboardHandler()

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

	routes.SetupUnprotectedRoutes(router, blogHandler, userHandler, dashboardHanlder)
	routes.SetupProtectedRoutes(router, pool, blogHandler, userHandler, dashboardHanlder)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
