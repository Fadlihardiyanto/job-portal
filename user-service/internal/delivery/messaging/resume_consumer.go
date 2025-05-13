package messaging

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type ResumeConsumer struct {
	Log *logrus.Logger
}

func NewResumeConsumer(log *logrus.Logger) *ResumeConsumer {
	return &ResumeConsumer{
		Log: log,
	}
}

func (c *ResumeConsumer) Consume(message *kafka.Message) error {
	// TODO process event
	c.Log.Infof("Received topic resumes with event: %s from partition %d", string(message.Value), message.TopicPartition.Partition)
	return nil
}
