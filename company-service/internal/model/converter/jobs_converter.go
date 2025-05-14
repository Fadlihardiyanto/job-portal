package converter

import (
	"company-service/internal/entity"
	"company-service/internal/model"
)

func JobsToResponse(jobs *entity.Jobs) *model.ResponseJobs {
	return &model.ResponseJobs{
		ID:            jobs.ID,
		JobsTitle:     jobs.JobsTitle,
		CompanyID:     jobs.CompanyID,
		Location:      jobs.Location,
		WorkspaceType: jobs.WorkspaceType,
		MinSalary:     jobs.MinSalary,
		MaxSalary:     jobs.MaxSalary,
		CreatedAt:     jobs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     jobs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func JobsToEvent(jobs *entity.Jobs) *model.JobsEvent {
	return &model.JobsEvent{
		ID:            jobs.ID,
		JobsTitle:     jobs.JobsTitle,
		CompanyID:     jobs.CompanyID,
		Location:      jobs.Location,
		WorkspaceType: jobs.WorkspaceType,
		MinSalary:     jobs.MinSalary,
		MaxSalary:     jobs.MaxSalary,
		CreatedAt:     jobs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     jobs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
