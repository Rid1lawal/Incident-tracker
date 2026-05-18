package middleware

import (
	"strconv"
	"time"

	"github.com/Reazy-ai/incident-tracker/internal/metrics"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		status := strconv.Itoa(c.Writer.Status())

		metrics.HttpRequestsTotal.WithLabelValues(
			path,
			c.Request.Method,
			status,
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			path,
			c.Request.Method,
		).Observe(duration)
	}
}
