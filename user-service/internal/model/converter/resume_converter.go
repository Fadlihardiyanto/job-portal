package converter

import (
	"strconv"
	"user-service/internal/entity"
	"user-service/internal/model"
)

func ResumeToResponse(resume *entity.Resume) *model.ResponseResume {
	return &model.ResponseResume{
		ID:         resume.ID,
		Name:       resume.Name,
		Attachment: resume.Attachment,
		UserID:     strconv.Itoa(resume.UserID),
		CreatedAt:  resume.CreatedAt,
		UpdatedAt:  resume.UpdatedAt,
	}
}
func ResumeToEvent(resume *entity.Resume) *model.ResumeEvent {
	return &model.ResumeEvent{
		ID:        resume.ID,
		Name:      resume.Name,
		UserID:    strconv.Itoa(resume.UserID),
		Status:    "queued",
		CreatedAt: resume.CreatedAt,
		UpdatedAt: resume.UpdatedAt,
	}
}
