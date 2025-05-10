package http

import (
	common "user-service/internal/common/error"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ResumeController struct {
	Log           *logrus.Logger
	ResumeUseCase *usecase.ResumeUseCase
}

func NewResumeController(log *logrus.Logger, resumeUseCase *usecase.ResumeUseCase) *ResumeController {
	return &ResumeController{
		Log:           log,
		ResumeUseCase: resumeUseCase,
	}
}

func (c *ResumeController) CreateResume(ctx *fiber.Ctx) error {
	request := new(model.RequestResume)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	response, err := c.ResumeUseCase.CreateResume(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create resume : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseResume]{
		Message: "Successfully created resume",
		Data:    *response,
	})
}
