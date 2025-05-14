package model

type ResponseJobs struct {
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

type RequestJobs struct {
	JobsTitle     string `json:"jobs_title" validate:"required"`
	CompanyID     string `json:"company_id" validate:"required"`
	Location      int    `json:"location" validate:"required"`
	WorkspaceType string `json:"workspace_type" validate:"required"`
	MinSalary     string `json:"min_salary" validate:"required"`
	MaxSalary     string `json:"max_salary" validate:"required"`
}

type RequestFindJobsByID struct {
	ID string `json:"id" validate:"required"`
}

type RequestFindJobsByCompanyID struct {
	CompanyID string `json:"company_id" validate:"required"`
}

type RequestUpdateJobs struct {
	ID            string `json:"id" validate:"omitempty"`
	JobsTitle     string `json:"jobs_title" validate:"omitempty"`
	CompanyID     string `json:"company_id" validate:"omitempty"`
	Location      int    `json:"location" validate:"omitempty"`
	WorkspaceType string `json:"workspace_type" validate:"omitempty"`
	MinSalary     string `json:"min_salary" validate:"omitempty"`
	MaxSalary     string `json:"max_salary" validate:"omitempty"`
	CreatedAt     string `json:"created_at" validate:"omitempty"`
	UpdatedAt     string `json:"updated_at" validate:"omitempty"`
}
