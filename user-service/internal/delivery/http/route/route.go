package route

import (
	"user-service/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	AuthController *http.AuthController
}

func (c *RouteConfig) Setup() {
	group := c.App.Group("/api/v1")

	// public route
	c.SetupGuestRoute(group)

	// auth middleware
	c.App.Use(c.AuthMiddleware)

	// auth route
	c.SetupAuthRoute(group)
}

func (c *RouteConfig) SetupGuestRoute(group fiber.Router) {
	// auth route
	authGroup := group.Group("/auth")
	authGroup.Post("/register", c.AuthController.Register)
}

func (c *RouteConfig) SetupAuthRoute(group fiber.Router) {
}
