package config

import (
	"company-service/internal/delivery/http"
	"company-service/internal/delivery/http/route"
	"company-service/internal/gateway/messaging"
	"company-service/internal/model"
	"company-service/internal/repository"
	"company-service/internal/usecase"

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

	// repo
	companyRepository := repository.NewCompanyRepository(config.Log)
	jobsRepository := repository.NewJobsRepository(config.Log)
	userJobsRepository := repository.NewUserJobsRepository(config.Log)

	// producer
	companyProducer := messaging.NewCompanyProducer(config.Producer, config.Log)
	jobsProducer := messaging.NewJobsProducer(config.Producer, config.Log)
	userJobsProducer := messaging.NewUserJobsProducer(config.Producer, config.Log)

	// usecase
	companyUseCase := usecase.NewCompanyUseCase(config.DB, config.Log, config.Validate, config.Viper, companyRepository, companyProducer)
	jobsUseCase := usecase.NewJobsUsecase(config.DB, config.Log, config.Validate, config.Viper, jobsRepository, jobsProducer)
	userJobsUseCase := usecase.NewUserJobsUsecase(config.DB, config.Log, config.Validate, config.Viper, userJobsRepository, jobsRepository, userJobsProducer)

	// controller
	companyController := http.NewCompanyController(config.Log, config.Viper, companyUseCase)
	jobsController := http.NewJobsController(config.Log, config.Viper, jobsUseCase)
	userJobsController := http.NewUserJobsController(config.Log, config.Viper, userJobsUseCase)

	routeConfig := route.RouteConfig{

		App:                config.App,
		CompanyController:  companyController,
		JobsController:     jobsController,
		UserJobsController: userJobsController,
	}

	routeConfig.SetupRoutes()
}
