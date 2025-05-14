package model

type ResponseCompany struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	City             string `json:"city"`
	OrganizationSize string `json:"organization_size"`
	Logo             string `json:"logo"`
	UserAccess       string `json:"user_access"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type RequestCompany struct {
	Name             string `json:"name" validate:"required"`
	City             string `json:"city" validate:"required"`
	OrganizationSize string `json:"organization_size" validate:"required"`
	Logo             string `json:"logo" validate:"required"`
	UserAccess       string `json:"user_access" validate:"required"`
}

type RequestFindCompanyByID struct {
	ID string `json:"id" validate:"required"`
}

type RequestUpdateCompany struct {
	ID               string `json:"id" validate:"omitempty"`
	Name             string `json:"name" validate:"omitempty"`
	City             string `json:"city" validate:"omitempty"`
	OrganizationSize string `json:"organization_size" validate:"omitempty"`
	Logo             string `json:"logo" validate:"omitempty"`
	UserAccess       string `json:"user_access" validate:"omitempty"`
	CreatedAt        string `json:"created_at" validate:"omitempty"`
	UpdatedAt        string `json:"updated_at" validate:"omitempty"`
}
