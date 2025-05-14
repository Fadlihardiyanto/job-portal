package messaging

import (
	"company-service/internal/model"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type JobsConsumer struct {
	Log *logrus.Logger
}

func NewJobsConsumer(log *logrus.Logger) *JobsConsumer {
	return &JobsConsumer{
		Log: log,
	}
}

func (c *JobsConsumer) Consume(message *kafka.Message) error {
	jobsEvent := new(model.JobsEvent)
	if err := json.Unmarshal(message.Value, jobsEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling Jobs event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic companies with event: %v from partition %d", jobsEvent, message.TopicPartition.Partition)
	return nil

}
