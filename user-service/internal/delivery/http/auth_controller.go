package http

import (
	"log"
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

	data, err := c.UsersUsecase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Message: "Successfully registered please check your email to verify your account",
		Data:    data,
	})
}

func (c *AuthController) Verify(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	request := &model.VerifyUserRequest{
		Token: token,
	}

	log.Printf("token: %s", token)

	response, err := c.UsersUsecase.Verify(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	log.Printf("response: %v", response)

	if response {
		return ctx.JSON(model.WebResponse[*string]{
			Message: "Successfully verified user",
		})
	}

	return ctx.JSON(model.WebResponse[*string]{
		Message: "Failed to verify user",
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UsersUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return common.HandleErrorResponse(ctx, err)
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
