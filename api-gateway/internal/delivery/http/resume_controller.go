package http

import (
	"api-gateway/internal/delivery/http/middleware"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ResumeController struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewResumeController(viper *viper.Viper, log *logrus.Logger) *ResumeController {
	return &ResumeController{
		Viper: viper,
		Log:   log,
	}
}

func (a *ResumeController) GetAllResume(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/resumes")
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

func (a *ResumeController) FindResumeByID(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/resumes/" + c.Params("id"))
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

func (a *ResumeController) GetResumeByUserID(c *fiber.Ctx) error {
	client := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(client)

	req := client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(a.Viper.GetString("USER_SERVICE_URL") + "/api/v1/resumes/user/" + c.Params("id"))
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

func (a *ResumeController) CreateResume(c *fiber.Ctx) error {
	client := &http.Client{}
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)

	auth := middleware.GetUser(c)

	log.Println("User ID: ", auth.ID)

	_ = writer.WriteField("user_id", strconv.Itoa(auth.ID))

	file, err := c.FormFile("attachment")
	if err != nil {
		a.Log.Error("Failed to get file from form: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	filePart, err := writer.CreateFormFile("attachment", file.Filename)
	if err != nil {
		a.Log.Error("Failed to create form file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	src, err := file.Open()
	if err != nil {
		a.Log.Error("Failed to open file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	defer src.Close()

	_, err = io.Copy(filePart, src)
	if err != nil {
		a.Log.Error("Failed to copy file data: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	err = writer.Close()
	if err != nil {
		a.Log.Error("Failed to close multipart writer: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	req, err := http.NewRequest("POST", a.Viper.GetString("USER_SERVICE_URL")+"/api/v1/resumes", reqBody)
	if err != nil {
		a.Log.Error("Failed to create new request: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		a.Log.Error("Failed to send request to user-service: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.Log.Error("Failed to read response body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(resp.StatusCode).Send(body)
}
