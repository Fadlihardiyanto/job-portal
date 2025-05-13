package messaging

import (
	"encoding/json"
	"log"
	"strconv"
	"user-service/internal/entity"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ResumeConsumer struct {
	Log              *logrus.Logger
	ResumeRepository *repository.ResumeRepository
	DB               *gorm.DB
}

func NewResumeConsumer(log *logrus.Logger, resumeRepository *repository.ResumeRepository, db *gorm.DB) *ResumeConsumer {
	return &ResumeConsumer{
		Log:              log,
		ResumeRepository: resumeRepository,
		DB:               db,
	}
}

func (c *ResumeConsumer) Consume(message *kafka.Message) error {
	var event model.ResumeEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		c.Log.Errorf("Failed to unmarshal resume event: %v", err)
		return err
	}

	c.Log.Infof("Processing resume event: %+v", event)

	userID, err := strconv.Atoi(event.UserID)
	if err != nil {
		c.Log.Errorf("Failed to convert UserID to int: %v", err)
		return err
	}

	resume := entity.Resume{
		Name:       event.Name,
		Attachment: event.Attachment,
		UserID:     userID,
	}

	log.Println("Resume event:", event)
	log.Println("Resume entity:", resume)

	if err := c.ResumeRepository.Create(c.DB, &resume); err != nil {
		c.Log.Errorf("Failed to save resume: %v", err)
	}

	c.Log.Infof("Resume saved successfully: %+v", resume)
	return nil
}
