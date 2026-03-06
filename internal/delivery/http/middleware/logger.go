// Package middleware provides Gin middleware functions for the HTTP delivery layer.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a Gin middleware that logs each request using the provided zap.Logger.
// It records: HTTP method, request path, client IP, status code, and request latency.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process the request first.
		c.Next()

		// Build the full path with query params if present.
		if query != "" {
			path = path + "?" + query
		}

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Log at different levels depending on the status code.
		switch {
		case statusCode >= 500:
			log.Error("server error", logFields...)
		case statusCode >= 400:
			log.Warn("client error", logFields...)
		default:
			log.Info("request", logFields...)
		}
	}
}

// Recovery returns a Gin middleware that recovers from panics, logs them,
// and returns a 500 Internal Server Error response.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   "internal server error",
				})
			}
		}()
		c.Next()
	}
}
