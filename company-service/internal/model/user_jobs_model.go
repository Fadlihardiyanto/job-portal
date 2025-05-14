package model

type UserJobsResponse struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	ResumeID  string `json:"resume_id"`
	JobID     string `json:"job_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserJobsRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	ResumeID  string `json:"resume_id" validate:"required"`
	JobID     string `json:"job_id" validate:"required"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
	UpdatedAt string `json:"updated_at" validate:"omitempty"`
}

type UserJobsUpdateRequest struct {
	ID        int    `json:"id" validate:"required"`
	UserID    int    `json:"user_id" validate:"omitempty"`
	ResumeID  int    `json:"resume_id" validate:"omitempty"`
	JobID     int    `json:"job_id" validate:"omitempty"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
	UpdatedAt string `json:"updated_at" validate:"omitempty"`
}
