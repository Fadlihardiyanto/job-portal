package config

import (
	"user-service/internal/delivery/http"
	"user-service/internal/delivery/http/route"
	"user-service/internal/gateway/messaging"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/usecase"

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

	// setup repositories
	userRepository := repository.NewUsersRepository(config.Log)
	resumeRepository := repository.NewResumeRepository(config.Log)

	// setup producers
	userProducer := messaging.NewUserProducer(config.Producer, config.Log)
	notifcationProducer := messaging.NewNotificationProducer(config.Producer, config.Log)
	resumeProducer := messaging.NewResumeProducer(config.Producer, config.Log)

	// setup use cases
	tokenUseCase := usecase.NewTokenUseCase(config.JWTConfig, config.Log)
	userUseCase := usecase.NewUsersUsecase(config.DB, config.Log, config.Validate, config.Viper, tokenUseCase, userRepository, userProducer, notifcationProducer)
	resumeUseCase := usecase.NewResumeUseCase(config.DB, config.Log, config.Validate, config.Viper, resumeRepository, resumeProducer)

	// setup controllers
	authController := http.NewAuthController(config.Log, userUseCase)
	userController := http.NewUserController(config.Log, userUseCase)
	resumeController := http.NewResumeController(config.Log, resumeUseCase, userUseCase)

	routeConfig := route.RouteConfig{

		App:              config.App,
		AuthController:   authController,
		UserController:   userController,
		ResumeController: resumeController,
	}

	routeConfig.SetupRoutes()
}
