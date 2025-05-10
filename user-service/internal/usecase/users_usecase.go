package usecase

import (
	"context"
	"fmt"
	"reflect"
	common "user-service/internal/common/error"
	commonUtil "user-service/internal/common/util"
	"user-service/internal/entity"
	"user-service/internal/gateway/messaging"
	"user-service/internal/model"
	"user-service/internal/model/converter"
	"user-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersUsecase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	Viper                *viper.Viper
	TokenUseCase         *TokenUseCase
	UsersRepository      *repository.UsersRepository
	UserProducer         *messaging.UserProducer
	NotificationProducer *messaging.NotificationProducer
}

func NewUsersUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, tokenUseCase *TokenUseCase, usersRepository *repository.UsersRepository, userProducer *messaging.UserProducer, notificationProducer *messaging.NotificationProducer) *UsersUsecase {
	return &UsersUsecase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		Viper:                viper,
		TokenUseCase:         tokenUseCase,
		UsersRepository:      usersRepository,
		UserProducer:         userProducer,
		NotificationProducer: notificationProducer,
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

	// generate token
	token, err := commonUtil.GenerateToken(32)
	if err != nil {
		c.Log.Warnf("Failed generate token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Password:   string(password),
		Email:      request.Email,
		Role:       request.Role,
		EmailToken: token,
	}

	if err := c.UsersRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	event := converter.UserToEvent(user)
	c.Log.Info("Publishing user created event")
	if err = c.UserProducer.Send(event); err != nil {
		c.Log.Warnf("Failed publish user created event : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// link verification
	linkVerification := fmt.Sprintf(
		"%s://%s:%s/api/v1/auth/verify?token=%s",
		c.Viper.GetString("WEB_PROTOCOL"),
		c.Viper.GetString("WEB_HOST"),
		c.Viper.GetString("WEB_PORT_GATEWAY"),
		user.EmailToken,
	)

	// nanti akan dikirimkan email lewat kafka di sini
	notifEvent := &model.NotificationEvent{
		ID:         user.ID,
		Email:      user.Email,
		TemplateID: "registration",
		Type:       "registration",
		Data: map[string]interface{}{
			"verification_link": linkVerification,
		},
	}

	c.Log.Info("Publishing notification event")
	if err = c.NotificationProducer.Send(notifEvent); err != nil {
		c.Log.Warnf("Failed publish notification event : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil

}

func (c *UsersUsecase) Verify(ctx context.Context, request *model.VerifyUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	user, err := c.UsersRepository.FindByEmailToken(tx, request.Token)
	if err != nil {
		c.Log.Warnf("User not found by token: %+v", err)
		return false, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.EmailVerifiedAt.Valid {
		c.Log.Warnf("Email already verified with token: %s", request.Token)
		return false, fiber.NewError(fiber.StatusConflict, "Email already verified")
	}

	if err := c.UsersRepository.UpdateEmailVerifiedAt(tx, user.ID); err != nil {
		c.Log.Warnf("Failed to update email_verified_at: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	event := converter.UserToEvent(user)
	c.Log.Info("Publishing user created event")

	if err := c.UserProducer.Send(event); err != nil {
		c.Log.Warnf("Failed publish user created event : %+v", err)
		return false, fiber.ErrInternalServerError
	}
	c.Log.Infof("User with email %s verified successfully", user.Email)

	return true, nil
}

func (c *UsersUsecase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	user := new(entity.User)
	if err := c.UsersRepository.FindByEmailVerified(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by id in Login : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "email not found or not verified or password not match")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, `email not found or not verified or password not match`)
	}

	// Generate JWT token
	accessToken, AccessExpiry, err := c.TokenUseCase.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.Log.Warnf("Failed to generate token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.Token = accessToken

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	event := converter.UserToEvent(user)
	c.Log.Info("Publishing user created event")
	if err := c.UserProducer.Send(event); err != nil {
		c.Log.Warnf("Failed publish user created event : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToLoginResponse(user, accessToken, AccessExpiry), nil

}

func (c *UsersUsecase) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := c.UsersRepository.FindByID(tx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Warnf("Failed to get user by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err == gorm.ErrRecordNotFound {
		c.Log.Warnf("User not found with id: %s", id)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UsersUsecase) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	users, err := c.UsersRepository.FindAll(tx)
	if err != nil {
		c.Log.Warnf("Failed to get all users : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := make([]model.UserResponse, len(users))
	for i, user := range users {
		response[i] = *converter.UserToResponse(&user)
	}

	return response, nil
}

func (c *UsersUsecase) UpdateUser(ctx context.Context, id string, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	user, err := c.UsersRepository.FindByID(tx, id)
	if err != nil {
		c.Log.Warnf("Failed to get user by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if user == nil {
		c.Log.Warnf("User not found with id: %s", id)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if request.Email != "" && request.Email != user.Email {
		exist, err := c.UsersRepository.IsExist(tx, request.Email)
		if err != nil {
			c.Log.Errorf("Error checking if user exists: %v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError)
		}
		if exist {
			c.Log.Warnf("User already exists with email: %s", request.Email)
			return nil, fiber.NewError(fiber.StatusConflict, "User already exists")
		}
		user.Email = request.Email
	}

	// Update allowed fields
	if request.Role != "" {
		user.Role = request.Role
	}
	if request.About != "" {
		user.About = request.About
	}
	if request.Photo != "" {
		user.Photo = request.Photo
	}
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}

	if err := c.UsersRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Publish user updated event
	event := converter.UserToEvent(user)
	c.Log.Info("Publishing user updated event")
	if err := c.UserProducer.Send(event); err != nil {
		c.Log.Warnf("Failed to publish user updated event : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}
