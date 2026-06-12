package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	DailyRequestLimit        = 10_000
	MonthlyUploadLimit int64 = 512 * 1024 * 1024 // 512 MB
)

func RequestQuotaMiddleWare(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Exclude option requests
		if c.Request.Method == http.MethodOptions {
			c.Next() // Skip this middleware and move to the next handler
			return
		}

		now := time.Now().UTC()

		key := "quota:requests:" + now.Format("2006-01-02")

		count, err := redisClient.Incr(c, key).Result()
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "quota check failed"})
			c.Abort()
			return
		}

		if count == 1 {
			midnight := time.Date(
				now.Year(),
				now.Month(),
				now.Day()+1,
				0,
				0,
				0,
				0,
				time.UTC,
			)
			redisClient.ExpireAt(c, key, midnight)
		}

		if count > DailyRequestLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "quota exceed limit",
				"limit": DailyRequestLimit,
				"used":  count,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func UploadQuotaMiddleWare(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.ContentLength <= 0 {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")

		if !strings.HasPrefix(contentType, "multipart/") {
			c.Next()
			return
		}

		if len(c.Request.Header.Get("Content-Type")) > 0 &&
			c.Request.Header.Get("Content-Type")[:9] != "multipart" {
			c.Next()
			return
		}

		now := time.Now().UTC()

		monthKey := "quota:uploads:" + now.Format("2006-01")

		current, err := redisClient.Get(c, monthKey).Int64()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "quota check failed",
			})
			c.Abort()
			return
		}

		projected := current + c.Request.ContentLength
		if projected > MonthlyUploadLimit {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "upload quota exceeded",
				"limit":   MonthlyUploadLimit,
				"used":    projected,
			})
			c.Abort()
			return
		}

		c.Next()

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			_, err = redisClient.IncrBy(c, monthKey, c.Request.ContentLength).Result()
			if err == nil && current == 0 {
				expireAtMonthEnd(c, redisClient, monthKey)
			}
		}

	}
}

func expireAtMonthEnd(
	ctx context.Context,
	rdb *redis.Client,
	key string,
) {
	now := time.Now().UTC()

	nextMonth := time.Date(
		now.Year(),
		now.Month()+1,
		1,
		0,
		0,
		0,
		0,
		now.Location(),
	)

	rdb.ExpireAt(ctx, key, nextMonth)
}
