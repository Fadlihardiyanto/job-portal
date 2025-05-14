package model

type CompanyEvent struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	City             string `json:"city"`
	OrganizationSize string `json:"organization_size"`
	Logo             string `json:"logo"`
	UserAccess       string `json:"user_access"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

func (c *CompanyEvent) GetId() int {
	return c.ID
}

func (c *CompanyEvent) GetKey() string {
	return c.UserAccess
}
