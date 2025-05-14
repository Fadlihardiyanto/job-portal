package model

type JobsEvent struct {
	ID            int    `json:"id"`
	JobsTitle     string `json:"jobs_title"`
	CompanyID     string `json:"company_id"`
	Location      int    `json:"location"`
	WorkspaceType string `json:"workspace_type"`
	MinSalary     string `json:"min_salary"`
	MaxSalary     string `json:"max_salary"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func (j *JobsEvent) GetId() int {
	return j.ID
}

func (j *JobsEvent) GetKey() string {
	return j.CompanyID
}
