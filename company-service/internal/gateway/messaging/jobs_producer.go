package messaging

import (
	"company-service/internal/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type JobsProducer struct {
	Producer[*model.JobsEvent]
}

func NewJobsProducer(producer *kafka.Producer, log *logrus.Logger) *JobsProducer {
	return &JobsProducer{
		Producer: Producer[*model.JobsEvent]{
			Producer: producer,
			Topic:    "jobs",
			Log:      log,
		},
	}
}
