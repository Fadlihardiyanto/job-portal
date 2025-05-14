package messaging

import (
	"company-service/internal/model"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type CompanyConsumer struct {
	Log *logrus.Logger
}

func NewCompanyConsumer(log *logrus.Logger) *CompanyConsumer {
	return &CompanyConsumer{
		Log: log,
	}
}

func (c *CompanyConsumer) Consume(message *kafka.Message) error {
	companyEvent := new(model.CompanyEvent)
	if err := json.Unmarshal(message.Value, companyEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling Company event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic companies with event: %v from partition %d", companyEvent, message.TopicPartition.Partition)
	return nil
}
