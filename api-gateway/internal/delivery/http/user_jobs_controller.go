package http

import (
	"api-gateway/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserJobsController struct {
	Viper            *viper.Viper
	Log              *logrus.Logger
	ResumeController *ResumeController
}

func NewUserJobsController(viper *viper.Viper, log *logrus.Logger, resumeController *ResumeController) *UserJobsController {
	return &UserJobsController{
		Viper:            viper,
		Log:              log,
		ResumeController: resumeController,
	}
}

func (c *UserJobsController) CreateUserJobs(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	if auth == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err := c.ResumeController.GetResumeByUserID(ctx)
	if err != nil {
		c.Log.Warnf("Failed to get resume by user ID : %+v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)
	req := client.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/user-jobs")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	req.SetBody(ctx.Body())

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()

	return ctx.Status(resp.StatusCode()).Send(body)
}
