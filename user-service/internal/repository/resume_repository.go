package repository

import (
	"user-service/internal/entity"

	"github.com/sirupsen/logrus"
)

type ResumeRepository struct {
	Repository[entity.Resume]
	Log *logrus.Logger
}

func NewResumeRepository(log *logrus.Logger) *ResumeRepository {
	return &ResumeRepository{
		Log: log,
	}
}
