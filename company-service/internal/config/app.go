package config

import (
	"company-service/internal/delivery/http/route"
	"company-service/internal/model"

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

	routeConfig := route.RouteConfig{

		App: config.App,
		// AuthController:   authController,
		// UserController:   userController,
		// ResumeController: resumeController,
	}

	routeConfig.SetupRoutes()
}
