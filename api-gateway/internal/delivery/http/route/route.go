package route

import (
	"api-gateway/internal/delivery/http"
	"api-gateway/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	AuthController   *http.AuthController
	UserController   *http.UsersController
	ResumeController *http.ResumeController

	AuthMiddleware fiber.Handler
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
	auth := group.Group("/auth")
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/register", c.AuthController.Register)
	auth.Get("/verify", c.AuthController.Verify)

}

func (c *RouteConfig) SetupAuthRoute(group fiber.Router) {
	users := group.Group("/users")
	users.Get("/", middleware.RoleMiddleware("admin"), c.UserController.GetAllUser)
	users.Get("/:id", middleware.RoleMiddleware("admin"), c.UserController.GetUserByID)
	users.Put("/", c.UserController.UpdateUser)

	resume := group.Group("/resumes")
	resume.Get("/", c.ResumeController.GetAllResume)
	resume.Get("/:id", c.ResumeController.FindResumeByID)
	resume.Get("/user/:id", c.ResumeController.GetResumeByUserID)
	resume.Post("/", c.ResumeController.CreateResume)
}
