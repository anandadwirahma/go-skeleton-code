// Package router wires all Gin routes and applies global middleware.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yourusername/go-skeleton-code/internal/delivery/http/handler"
	"github.com/yourusername/go-skeleton-code/internal/delivery/http/middleware"
)

// Setup creates and configures the Gin engine with all application routes.
// It accepts the concrete handler structs and a logger for middleware.
func Setup(
	log *zap.Logger,
	contactHandler *handler.ContactHandler,
) *gin.Engine {
	engine := gin.New() // Use gin.New() instead of gin.Default() for custom middleware

	// --- Global Middleware ---
	engine.Use(middleware.Recovery(log))
	engine.Use(middleware.Logger(log))

	// --- Health check ---
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --- API v1 group ---
	v1 := engine.Group("/api/v1")
	{
		contacts := v1.Group("/contacts")
		{
			contacts.POST("", contactHandler.Create)
			contacts.GET("", contactHandler.GetAll)
			contacts.GET("/:id", contactHandler.GetByID)
			contacts.PUT("/:id", contactHandler.Update)
			contacts.DELETE("/:id", contactHandler.Delete)
		}
	}

	return engine
}
