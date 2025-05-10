package http

import (
	common "user-service/internal/common/error"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log          *logrus.Logger
	UsersUsecase *usecase.UsersUsecase
}

func NewAuthController(log *logrus.Logger, usersUsecase *usecase.UsersUsecase) *AuthController {
	return &AuthController{
		Log:          log,
		UsersUsecase: usersUsecase,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	_, err = c.UsersUsecase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[*string]{
		Message: "Successfully registered please check your email to verify your account",
	})
}
