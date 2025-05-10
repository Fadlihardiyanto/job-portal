package http

import (
	"api-gateway/internal/model"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthController struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewAuthController(viper *viper.Viper, log *logrus.Logger) *AuthController {
	return &AuthController{
		Viper: viper,
		Log:   log,
	}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/auth/login")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
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
	var wrapped model.WrappedUserResponse
	if err := json.Unmarshal(resp.Body(), &wrapped); err != nil {
		a.Log.Error("Failed to unmarshal wrapped response: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// set access token to cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    wrapped.Data.AccessToken,
		Expires:  wrapped.Data.AccessExpiry,
		HTTPOnly: true,
		Secure:   true,
	})

	body := resp.Body()

	return c.Status(resp.StatusCode()).Send(body)

}

func (a *AuthController) Register(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/auth/register")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
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

func (a *AuthController) Verify(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/auth/verify?token=" + c.Query("token"))
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
