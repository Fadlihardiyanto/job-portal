package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JobsController struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewJobsController(viper *viper.Viper, log *logrus.Logger) *JobsController {
	return &JobsController{
		Viper: viper,
		Log:   log,
	}
}

func (c *JobsController) GetAllJobs(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/jobs")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))
	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}

func (c *JobsController) GetJobsByID(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/jobs/" + ctx.Params("id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}

func (c *JobsController) GetJobsByCompanyID(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/jobs/company/" + ctx.Params("company_id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}

func (c *JobsController) CreateJob(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/jobs")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	req.SetBody(ctx.Body())

	if err := ctx.BodyParser(req.Body()); err != nil {
		c.Log.Error("Failed to parse request body: ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}

func (c *JobsController) UpdateJob(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPut)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/jobs/" + ctx.Params("id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	req.SetBody(ctx.Body())

	if err := ctx.BodyParser(req.Body()); err != nil {
		c.Log.Error("Failed to parse request body: ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to auth service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}
