package converter

import (
	"company-service/internal/entity"
	"company-service/internal/model"
	"strconv"
)

func ResumeToResponse(resume *entity.Resume) *model.ResponseResume {
	return &model.ResponseResume{
		ID:         resume.ID,
		Name:       resume.Name,
		Attachment: resume.Attachment,
		UserID:     strconv.Itoa(resume.UserID),
		CreatedAt:  resume.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  resume.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
func ResumeToEvent(resume *entity.Resume) *model.ResumeEvent {
	return &model.ResumeEvent{
		ID:        resume.ID,
		Name:      resume.Name,
		UserID:    strconv.Itoa(resume.UserID),
		Status:    "queued",
		CreatedAt: resume.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: resume.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func EventToResponse(event *model.ResumeEvent) *model.ResponseResume {
	return &model.ResponseResume{
		Name:       event.Name,
		UserID:     event.UserID,
		Attachment: event.Attachment,
		Status:     event.Status,
		CreatedAt:  event.CreatedAt,
		UpdatedAt:  event.UpdatedAt,
	}
}
