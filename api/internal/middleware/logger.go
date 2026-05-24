package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start)
		requestID, _ := c.Get(RequestIDKey)

		log.Info().
			Str("request_id", requestID.(string)).
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("incoming request")
	}
}
