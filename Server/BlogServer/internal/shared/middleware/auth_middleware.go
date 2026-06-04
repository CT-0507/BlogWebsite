package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	commondb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/repository/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleWare(pool *pgxpool.Pool, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := utils.GetAccessToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		tokenVer := claims.TokenVersion

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		key := fmt.Sprintf(
			"tokenCheck:user:%s",
			claims.UserID,
		)
		currentVer := -1
		val, err := redisClient.Get(ctx, key).Result()
		if err == nil {
			log.Println("Redis: AuthMiddleware: Hit: ")
			intV, err := strconv.Atoi(val)
			if err != nil {
				// Handle error (e.g., if string is "42a")
				log.Println(err)
			} else {
				currentVer = intV
			}
		} else if err != redis.Nil {
			log.Println("Redis: AuthMiddleware: cannot get: ", err)
		}

		// Avoid invalid token spam
		// if currentVer == -1 {
		// Guarantee correctness
		if currentVer == -1 || tokenVer != currentVer {
			commonQueries := commondb.New(pool)
			currentVerInt4, err := commonQueries.GetUserTokenVersionByID(ctx, userUUID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token DB"})
				c.Abort()
				return
			}
			currentVer = int(currentVerInt4.Int32)
			_ = redisClient.Set(ctx, key, currentVer, time.Hour).Err()
		}

		if tokenVer != currentVer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Compare"})
			c.Abort()
			return
		}

		c.Set("avatar", claims.Avatar)
		c.Set("username", claims.Username)
		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("role", claims.Roles[0])

		c.Next()

	}
}
