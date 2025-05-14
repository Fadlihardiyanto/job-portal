package converter

import (
	"company-service/internal/entity"
	"company-service/internal/model"
	"strconv"
)

func UserJobsToResponse(userJobs *entity.UserJobs) *model.UserJobsResponse {
	return &model.UserJobsResponse{
		ID:        userJobs.ID,
		UserID:    strconv.Itoa(userJobs.UserID),
		ResumeID:  strconv.Itoa(userJobs.ResumeID),
		JobID:     strconv.Itoa(userJobs.JobsID),
		CreatedAt: userJobs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: userJobs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func EventToResponseUserJobs(userJobs *model.UserJobsEvent) *model.UserJobsResponse {
	return &model.UserJobsResponse{
		ID:        userJobs.ID,
		UserID:    userJobs.UserID,
		ResumeID:  userJobs.ResumeID,
		JobID:     userJobs.JobID,
		CreatedAt: userJobs.CreatedAt,
		UpdatedAt: userJobs.UpdatedAt,
	}
}
