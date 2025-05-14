package messaging

import (
	"company-service/internal/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type UserJobsProducer struct {
	Producer[*model.UserJobsEvent]
}

func NewUserJobsProducer(producer *kafka.Producer, log *logrus.Logger) *UserJobsProducer {
	return &UserJobsProducer{
		Producer: Producer[*model.UserJobsEvent]{
			Producer: producer,
			Topic:    "user_jobs",
			Log:      log,
		},
	}
}
