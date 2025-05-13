package http

import (
	"time"
	common "user-service/internal/common/error"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ResumeController struct {
	Log           *logrus.Logger
	ResumeUseCase *usecase.ResumeUseCase
	UsersUseCase  *usecase.UsersUsecase
}

func NewResumeController(log *logrus.Logger, resumeUseCase *usecase.ResumeUseCase, usersUseCase *usecase.UsersUsecase) *ResumeController {
	return &ResumeController{
		Log:           log,
		ResumeUseCase: resumeUseCase,
		UsersUseCase:  usersUseCase,
	}
}

func (c *ResumeController) GetAllResumes(ctx *fiber.Ctx) error {
	resumes, err := c.ResumeUseCase.GetAllResume(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to get all resumes : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[[]model.ResponseResume]{
		Message: "Successfully retrieved all resumes",
		Data:    resumes,
	})
}

func (c *ResumeController) GetByUserID(ctx *fiber.Ctx) error {
	request := &model.RequestFindResumeByUser{
		UserID: ctx.Params("user_id"),
	}

	resumes, err := c.ResumeUseCase.GetByUserID(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get resumes by user ID : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[[]model.ResponseResume]{
		Message: "Successfully retrieved resumes",
		Data:    resumes,
	})
}

func (c *ResumeController) FindByID(ctx *fiber.Ctx) error {
	request := &model.RequestFindResumeByID{
		ID: ctx.Params("id"),
	}

	resume, err := c.ResumeUseCase.FindByID(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to find resume by ID : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseResume]{
		Message: "Successfully retrieved resume",
		Data:    *resume,
	})
}

func (c *ResumeController) CreateResume(ctx *fiber.Ctx) error {
	userID := ctx.FormValue("user_id")
	if userID == "" {
		return common.HandleErrorResponse(ctx, fiber.NewError(fiber.StatusBadRequest, "user_id is required"))
	}

	fileHeader, err := ctx.FormFile("attachment")
	if err != nil {
		c.Log.Warnf("Failed to get file from form: %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	nameFile := fileHeader.Filename
	timestamp := time.Now().Unix()
	uniqueID := uuid.NewString()[:8]
	storedFileName := fmt.Sprintf("%d_%s_%s", timestamp, uniqueID, nameFile)
	savePath := "uploads/resumes/" + storedFileName
	if err := ctx.SaveFile(fileHeader, savePath); err != nil {
		c.Log.Warnf("Failed to save file: %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	request := &model.RequestResume{
		Name:       nameFile,
		Attachment: savePath,
		UserID:     userID,
	}

	_, err = c.UsersUseCase.GetUserByID(ctx.UserContext(), request.UserID)
	if err != nil {
		c.Log.Warnf("Failed to find user by ID : %+v", err)
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
