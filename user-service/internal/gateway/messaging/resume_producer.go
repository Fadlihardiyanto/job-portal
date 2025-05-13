package messaging

import (
	"user-service/internal/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type ResumeProducer struct {
	Producer[*model.ResumeEvent]
}

func NewResumeProducer(producer *kafka.Producer, log *logrus.Logger) *ResumeProducer {
	return &ResumeProducer{
		Producer: Producer[*model.ResumeEvent]{
			Producer: producer,
			Topic:    "resumes",
			Log:      log,
		},
	}
}
