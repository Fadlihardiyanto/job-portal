package config

import (
	"user-service/internal/delivery/http"
	"user-service/internal/delivery/http/route"
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

	// setup use cases
	userUseCase := usecase.NewUsersUsecase(config.DB, config.Log, config.Validate, config.Viper, userRepository)
	// setup controllers
	authController := http.NewAuthController(config.Log, userUseCase)

	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
	}

	routeConfig.Setup()
}
