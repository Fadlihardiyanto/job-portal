package entity

import "time"

type Jobs struct {
	ID            int       `gorm:"column:id;primaryKey"`
	JobsTitle     string    `gorm:"column:jobs_title"`
	CompanyID     string    `gorm:"column:company_id"`
	Location      int       `gorm:"column:location"`
	WorkspaceType string    `gorm:"column:workspace_type"`
	MinSalary     string    `gorm:"column:min_salary"`
	MaxSalary     string    `gorm:"column:max_salary"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	Company       Company   `gorm:"foreignKey:CompanyID;references:ID"`
}
