package repository

import (
	"company-service/internal/entity"

	"github.com/sirupsen/logrus"
)

type UserJobsRepository struct {
	Repository[entity.UserJobs]
	Log *logrus.Logger
}

func NewUserJobsRepository(log *logrus.Logger) *UserJobsRepository {
	return &UserJobsRepository{
		Log: log,
	}
}
