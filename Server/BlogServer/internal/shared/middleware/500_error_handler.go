package middleware

import (
	"errors"
	"log"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		log.Printf("[ERROR] %v\n", err)

		var errDup *contracts.ErrDuplicate
		if errors.As(err, &errDup) {
			c.AbortWithStatusJSON(409, gin.H{
				"error": "Resource already exists",
			})
			return
		}

		c.AbortWithStatusJSON(500, gin.H{
			"error": "Internal Server Error",
		})
	}
}
