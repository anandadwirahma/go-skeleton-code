package route

import (
	"github.com/gofiber/fiber/v3"

	"github.com/yourusername/go-skeleton-code/internal/delivery/http"
)

type RouteConfig struct {
	App               *fiber.App
	ExampleController *http.ExampleController
}

func (c *RouteConfig) Setup() {
	c.SetupRoute()
}

func (c *RouteConfig) SetupRoute() {
	c.App.Post("/api/example", c.ExampleController.Create)
}
