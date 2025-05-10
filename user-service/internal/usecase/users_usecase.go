package usecase

import (
	"context"
	"reflect"
	common "user-service/internal/common/error"
	commonUtil "user-service/internal/common/util"
	"user-service/internal/entity"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersUsecase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	Viper           *viper.Viper
	UsersRepository *repository.UsersRepository
}

func NewUsersUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper) *UsersUsecase {
	return &UsersUsecase{
		DB:       db,
		Log:      log,
		Validate: validate,
		Viper:    viper,
	}
}

func (c *UsersUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	exist, err := c.UsersRepository.IsExist(tx, request.Email)
	if err != nil {
		c.Log.Errorf("Error checking if user exists: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}

	if exist {
		c.Log.Warnf("User already exists with email: %s", request.Email)
		return nil, fiber.NewError(fiber.StatusConflict, "User already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.Users{
		ID:       request.ID,
		Password: string(password),
		Email:    request.Email,
		// EmailToken: token,
	}

	if err := c.UsersRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return nil, nil
}
