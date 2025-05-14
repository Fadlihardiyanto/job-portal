package messaging

import (
	"company-service/internal/entity"
	"company-service/internal/model"
	"company-service/internal/repository"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserJobsConsumer struct {
	Log          *logrus.Logger
	UserJobsRepo *repository.UserJobsRepository
	DB           *gorm.DB
}

func NewUserJobsConsumer(log *logrus.Logger, userJobsRepo *repository.UserJobsRepository, db *gorm.DB) *UserJobsConsumer {
	return &UserJobsConsumer{
		Log:          log,
		UserJobsRepo: userJobsRepo,
		DB:           db,
	}
}

func (c *UserJobsConsumer) Consume(message *kafka.Message) error {
	event := new(model.UserJobsEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		c.Log.WithError(err).Error("error unmarshalling UserJobs event")
		return err
	}

	userID, _ := strconv.Atoi(event.UserID)
	resumeID, _ := strconv.Atoi(event.ResumeID)
	jobsID, _ := strconv.Atoi(event.JobID)
	createdAt, _ := time.Parse(time.RFC3339, event.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, event.UpdatedAt)

	userJobs := &entity.UserJobs{
		UserID:    userID,
		JobsID:    jobsID,
		ResumeID:  resumeID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	log.Println("UserJobs entity:", userJobs)

	if err := c.UserJobsRepo.Create(c.DB, userJobs); err != nil {
		c.Log.Errorf("Failed to save user jobs: %v", err)
	}

	c.Log.Infof("Received topic user_jobs with event: %s from partition %d", string(message.Value), message.TopicPartition.Partition)
	// TODO process event
	return nil
}
