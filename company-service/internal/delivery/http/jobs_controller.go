package http

import (
	"company-service/internal/model"
	"company-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JobsController struct {
	Log         *logrus.Logger
	Viper       *viper.Viper
	JobsUsecase *usecase.JobsUsecase
}

func NewJobsController(log *logrus.Logger, viper *viper.Viper, jobsUsecase *usecase.JobsUsecase) *JobsController {
	return &JobsController{
		Log:         log,
		Viper:       viper,
		JobsUsecase: jobsUsecase,
	}
}

func (c *JobsController) CreateJob(ctx *fiber.Ctx) error {
	request := new(model.RequestJobs)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	job, err := c.JobsUsecase.CreateJob(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create job : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create job")
	}

	return ctx.JSON(model.WebResponse[model.ResponseJobs]{
		Message: "Successfully created job",
		Data:    *job,
	})

}

func (c *JobsController) GetAllJobs(ctx *fiber.Ctx) error {
	jobs, err := c.JobsUsecase.GetAllJobs(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to get all jobs : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get all jobs")
	}

	return ctx.JSON(model.WebResponse[[]model.ResponseJobs]{
		Message: "Successfully retrieved all jobs",
		Data:    jobs,
	})
}

func (c *JobsController) GetJobsByCompanyID(ctx *fiber.Ctx) error {
	request := &model.RequestFindJobsByCompanyID{
		CompanyID: ctx.Params("company_id"),
	}

	jobs, err := c.JobsUsecase.GetJobsByCompanyID(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get jobs by company ID : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get jobs by company ID")
	}

	return ctx.JSON(model.WebResponse[[]model.ResponseJobs]{
		Message: "Successfully retrieved jobs",
		Data:    jobs,
	})
}

func (c *JobsController) FindJobByID(ctx *fiber.Ctx) error {
	request := &model.RequestFindJobsByID{
		ID: ctx.Params("id"),
	}

	job, err := c.JobsUsecase.FindByID(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get job by ID : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get job by ID")
	}

	return ctx.JSON(model.WebResponse[model.ResponseJobs]{
		Message: "Successfully retrieved job",
		Data:    *job,
	})
}

func (c *JobsController) UpdateJob(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateJobs)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	job, err := c.JobsUsecase.UpdateJob(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update job : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update job")
	}

	return ctx.JSON(model.WebResponse[model.ResponseJobs]{
		Message: "Successfully updated job",
		Data:    *job,
	})
}
