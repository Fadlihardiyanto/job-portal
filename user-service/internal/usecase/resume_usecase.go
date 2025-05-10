package usecase

import (
	"context"
	"reflect"
	common "user-service/internal/common/error"
	commonUtil "user-service/internal/common/util"
	"user-service/internal/entity"
	"user-service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ResumeUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Viper    *viper.Viper
}

func NewResumeUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper) *ResumeUseCase {
	return &ResumeUseCase{
		DB:       db,
		Log:      log,
		Validate: validate,
		Viper:    viper,
	}
}

func (c *ResumeUseCase) CreateResume(ctx context.Context, request *model.RequestResume) (*model.ResponseResume, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	resume := entity.Resume{
		Name:      request.Name,
		Attacment: request.Attachment,
		UserID:    request.UserID,
	}

	if err := tx.Create(&resume).Error; err != nil {
		c.Log.Errorf("Failed to create resume: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create resume")
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}
	response := &model.ResponseResume{
		ID:         resume.ID,
		Name:       resume.Name,
		Attachment: resume.Attacment,
		UserID:     resume.UserID,
		CreatedAt:  resume.CreatedAt,
		UpdatedAt:  resume.UpdatedAt,
		Users: model.User{
			ID:        resume.Users.ID,
			FirstName: resume.Users.FirstName,
			LastName:  resume.Users.LastName,
			Email:     resume.Users.Email,
			About:     resume.Users.About,
			Photo:     resume.Users.Photo,
			Role:      resume.Users.Role,
			IsActive:  resume.Users.IsActive,
		},
	}

	return response, nil

}
