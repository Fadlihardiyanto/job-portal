package route

import "github.com/gofiber/fiber/v2"

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) SetupRoutes() {

}
