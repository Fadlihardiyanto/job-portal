package route

import (
	"user-service/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	AuthMiddleware   fiber.Handler
	AuthController   *http.AuthController
	UserController   *http.UserController
	ResumeController *http.ResumeController
}

func (c *RouteConfig) SetupRoutes() {
	api := c.App.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Get("/verify", c.AuthController.Verify)

	users := api.Group("/users")
	users.Get("/", c.UserController.GetAllUsers)
	users.Get("/:id", c.UserController.GetUserByID)
	users.Put("/:id", c.UserController.UpdateUser)

	// Resume routes
	resumes := api.Group("/resumes")
	resumes.Post("/", c.ResumeController.CreateResume)
}
