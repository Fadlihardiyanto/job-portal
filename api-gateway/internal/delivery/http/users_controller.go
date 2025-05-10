package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UsersController struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewUsersController(viper *viper.Viper, log *logrus.Logger) *UsersController {
	return &UsersController{
		Viper: viper,
		Log:   log,
	}
}

func (a *UsersController) GetAllUser(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/users")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", c.Get("Authorization"))

	if err := client.Parse(); err != nil {
		a.Log.Error("Failed to parse HTTP client: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		a.Log.Error("Failed to send request to auth service: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	body := resp.Body()

	return c.Status(resp.StatusCode()).Send(body)

}

func (a *UsersController) GetUserByID(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/users/" + c.Params("id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", c.Get("Authorization"))

	if err := client.Parse(); err != nil {
		a.Log.Error("Failed to parse HTTP client: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		a.Log.Error("Failed to send request to auth service: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()

	return c.Status(resp.StatusCode()).Send(body)

}

func (a *UsersController) UpdateUser(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPut)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/users/" + c.Params("id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", c.Get("Authorization"))
	req.SetBody(c.Body())

	if err := client.Parse(); err != nil {
		a.Log.Error("Failed to parse HTTP client: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		a.Log.Error("Failed to send request to auth service: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	body := resp.Body()

	return c.Status(resp.StatusCode()).Send(body)
}
