package http

import (
	common "company-service/internal/common/error"
	"company-service/internal/model"
	"company-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserJobsController struct {
	Log             *logrus.Logger
	Viper           *viper.Viper
	UserJobsUseCase *usecase.UserJobsUsecase
}

func NewUserJobsController(log *logrus.Logger, viper *viper.Viper, userJobsUseCase *usecase.UserJobsUsecase) *UserJobsController {
	return &UserJobsController{
		Log:             log,
		Viper:           viper,
		UserJobsUseCase: userJobsUseCase,
	}
}

func (c *UserJobsController) CreateUserJob(ctx *fiber.Ctx) error {
	request := new(model.UserJobsRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	userJob, err := c.UserJobsUseCase.CreateUserJobs(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create user job : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.UserJobsResponse]{
		Message: "Successfully created user job",
		Data:    *userJob,
	})
}
