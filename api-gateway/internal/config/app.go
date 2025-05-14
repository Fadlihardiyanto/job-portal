package config

import (
	"api-gateway/internal/delivery/http"
	"api-gateway/internal/delivery/http/middleware"
	"api-gateway/internal/delivery/http/route"
	"api-gateway/internal/model"
	"api-gateway/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Log       *logrus.Logger
	Validate  *validator.Validate
	Viper     *viper.Viper
	Producer  *kafka.Producer
	JWTConfig *model.JWTConfig
}

func Bootstrap(config *BootstrapConfig) {

	// usecase
	tokenUseCase := usecase.NewTokenUseCase(config.JWTConfig, config.Log)

	// setup controllers
	authController := http.NewAuthController(config.Viper, config.Log)
	userController := http.NewUsersController(config.Viper, config.Log)
	resumeController := http.NewResumeController(config.Viper, config.Log)
	companyController := http.NewCompanyController(config.Viper, config.Log)
	jobsController := http.NewJobsController(config.Viper, config.Log)
	userJobsController := http.NewUserJobsController(config.Viper, config.Log, resumeController)

	// setup middleware
	authMiddleware := middleware.NewAuth(tokenUseCase)

	routeConfig := route.RouteConfig{

		App:                config.App,
		AuthController:     authController,
		AuthMiddleware:     authMiddleware,
		UserController:     userController,
		ResumeController:   resumeController,
		CompanyController:  companyController,
		JobsController:     jobsController,
		UserJobsController: userJobsController,
	}

	routeConfig.Setup()
}
