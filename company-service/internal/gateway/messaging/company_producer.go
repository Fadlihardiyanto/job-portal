package messaging

import (
	"company-service/internal/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type CompanyProducer struct {
	Producer[*model.CompanyEvent]
}

func NewCompanyProducer(producer *kafka.Producer, log *logrus.Logger) *CompanyProducer {
	return &CompanyProducer{
		Producer: Producer[*model.CompanyEvent]{
			Producer: producer,
			Topic:    "companies",
			Log:      log,
		},
	}
}
