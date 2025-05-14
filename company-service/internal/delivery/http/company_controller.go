package http

import (
	common "company-service/internal/common/error"
	"company-service/internal/model"
	"company-service/internal/usecase"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CompanyController struct {
	Log            *logrus.Logger
	Viper          *viper.Viper
	CompanyUseCase *usecase.CompanyUseCase
}

func NewCompanyController(log *logrus.Logger, viper *viper.Viper, companyUseCase *usecase.CompanyUseCase) *CompanyController {
	return &CompanyController{
		Log:            log,
		Viper:          viper,
		CompanyUseCase: companyUseCase,
	}
}

func (c *CompanyController) CreateCompany(ctx *fiber.Ctx) error {
	request := new(model.RequestCompany)
	log.Println("Request body: ", request)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	company, err := c.CompanyUseCase.CreateCompany(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create company : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseCompany]{
		Message: "Successfully created company",
		Data:    *company,
	})
}

func (c *CompanyController) GetAllCompanies(ctx *fiber.Ctx) error {
	companies, err := c.CompanyUseCase.GetAllCompany(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to get all companies : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[[]model.ResponseCompany]{
		Message: "Successfully retrieved all companies",
		Data:    companies,
	})
}

func (c *CompanyController) FindByID(ctx *fiber.Ctx) error {
	request := &model.RequestFindCompanyByID{
		ID: ctx.Params("id"),
	}

	company, err := c.CompanyUseCase.FindByID(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get company by ID : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseCompany]{
		Message: "Successfully retrieved company",
		Data:    *company,
	})
}

func (c *CompanyController) UpdateCompany(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateCompany)

	log.Println("Request body: ", request)
	log.Println("Request params: ", ctx.Body())

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	if err := ctx.ParamsParser(request); err != nil {
		c.Log.Warnf("Failed to parse request params : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	company, err := c.CompanyUseCase.UpdateCompany(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update company : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseCompany]{
		Message: "Successfully updated company",
		Data:    *company,
	})
}

func (c *CompanyController) UpdateCompanyByIDAndUserAccess(ctx *fiber.Ctx) error {
	log.Println("Request body)")

	request := &model.RequestUpdateCompanyByIDAndUserAccess{
		ID:         ctx.Params("id"),
		UserAccess: ctx.Params("user_access"),
	}
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	if err := ctx.ParamsParser(request); err != nil {
		c.Log.Warnf("Failed to parse request params : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	company, err := c.CompanyUseCase.UpdateCompanyByIDAndUserAccess(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update company by ID and user access : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.ResponseCompany]{
		Message: "Successfully updated company by ID and user access",
		Data:    *company,
	})
}
