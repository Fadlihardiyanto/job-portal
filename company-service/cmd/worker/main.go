package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"company-service/internal/config"
	"company-service/internal/delivery/messaging"
	"company-service/internal/repository"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)

	logger.Info("Starting worker service")

	ctx, cancel := context.WithCancel(context.Background())

	go RunCompanyConsumer(logger, viperConfig, ctx)
	go RunJobsConsumer(logger, viperConfig, ctx)
	go RunUserJobsConsumer(logger, viperConfig, ctx, db)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	stop := false
	for !stop {
		s := <-terminateSignals
		logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
		cancel()
		stop = true
	}

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

func RunCompanyConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("setup company consumer")
	companyConsumer := config.NewKafkaConsumer(viperConfig, logger)
	companyHandler := messaging.NewCompanyConsumer(logger)
	messaging.ConsumeTopic(ctx, companyConsumer, "companies", logger, companyHandler.Consume)
}

func RunJobsConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("setup jobs consumer")
	jobsConsumer := config.NewKafkaConsumer(viperConfig, logger)
	jobsHandler := messaging.NewJobsConsumer(logger)
	messaging.ConsumeTopic(ctx, jobsConsumer, "jobs", logger, jobsHandler.Consume)
}

func RunUserJobsConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context, db *gorm.DB) {
	logger.Info("setup user jobs consumer")
	userJobsConsumer := config.NewKafkaConsumer(viperConfig, logger)
	userJobsRepository := repository.NewUserJobsRepository(logger)
	userJobsHandler := messaging.NewUserJobsConsumer(logger, userJobsRepository, db)
	messaging.ConsumeTopic(ctx, userJobsConsumer, "user_jobs", logger, userJobsHandler.Consume)
}
