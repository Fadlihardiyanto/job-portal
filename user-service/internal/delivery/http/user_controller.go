package http

import (
	common "user-service/internal/common/error"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UsersUsecase
}

func NewUserController(log *logrus.Logger, useCase *usecase.UsersUsecase) *UserController {
	return &UserController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.UseCase.GetAllUsers(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to get all users : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{
		Message: "Successfully retrieved all users",
		Data:    users,
	})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user ID is required")
	}

	user, err := c.UseCase.GetUserByID(ctx.UserContext(), userID)
	if err != nil {
		c.Log.Warnf("Failed to get user by ID : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.UserResponse]{
		Message: "Successfully retrieved user",
		Data:    *user,
	})
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user ID is required")
	}

	request := new(model.UpdateUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	user, err := c.UseCase.UpdateUser(ctx.UserContext(), userID, request)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[model.UserResponse]{
		Message: "Successfully updated user",
		Data:    *user,
	})
}
