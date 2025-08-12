package infrastructure

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AIRateLimitMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Minute), 5) 

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests to AI service. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
