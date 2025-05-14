package model

type UserJobsEvent struct {
	ID        int    `json:"id"`
	JobID     string `json:"job_id"`
	UserID    string `json:"user_id"`
	ResumeID  string `json:"resume_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *UserJobsEvent) GetId() int {
	return u.ID
}

func (u *UserJobsEvent) GetKey() string {
	return u.UserID
}
