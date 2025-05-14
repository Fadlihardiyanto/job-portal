package http

import (
	"api-gateway/internal/delivery/http/middleware"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CompanyController struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewCompanyController(viper *viper.Viper, log *logrus.Logger) *CompanyController {
	return &CompanyController{
		Viper: viper,
		Log:   log,
	}
}

func (c *CompanyController) GetAllCompany(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/company")
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

func (c *CompanyController) GetCompanyByID(ctx *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/company/" + ctx.Params("id"))
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

func (c *CompanyController) CreateCompany(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userID := auth.ID
	if userID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	var companyData map[string]interface{}
	if err := json.Unmarshal(ctx.Body(), &companyData); err != nil {
		c.Log.Error("Failed to parse request body: ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	companyData["user_access"] = strconv.Itoa(userID)

	updatedBody, err := json.Marshal(companyData)
	if err != nil {
		c.Log.Error("Failed to marshal updated JSON: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/company")
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.SetBody(updatedBody)

	if err := client.Parse(); err != nil {
		c.Log.Error("Failed to parse HTTP client: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		c.Log.Error("Failed to send request to company service: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	body := resp.Body()
	return ctx.Status(resp.StatusCode()).Send(body)
}

func (c *CompanyController) UpdateCompany(ctx *fiber.Ctx) error {

	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPut)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/company/" + ctx.Params("id"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	req.SetBody(ctx.Body())

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

func (c *CompanyController) UpdateCompanyByIDAndAccess(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	userID := auth.ID

	log.Println("User ID: ", userID)

	if userID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	if strconv.Itoa(userID) != ctx.Params("user_access") {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to access this resource",
		})
	}

	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodPut)
	req.SetRequestURI(c.Viper.GetString("COMPANY_SERVICE_URL") + "/api/v1/company/" + ctx.Params("id") + "/access/" + ctx.Params("user_access"))
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.Set("Authorization", ctx.Get("Authorization"))

	req.SetBody(ctx.Body())

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
