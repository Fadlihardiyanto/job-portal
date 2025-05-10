package entity

type Jobs struct {
	ID            int     `gorm:"column:id;primaryKey"`
	JobsTitle     string  `gorm:"column:jobs_title"`
	CompanyID     int     `gorm:"column:company_id"`
	Location      string  `gorm:"column:location"`
	WorkspaceType string  `gorm:"column:workspace_type"`
	MinSalary     string  `gorm:"column:min_salary"`
	MaxSalary     string  `gorm:"column:max_salary"`
	CreatedAt     string  `gorm:"column:created_at"`
	UpdatedAt     string  `gorm:"column:updated_at"`
	Company       Company `gorm:"foreignKey:CompanyID;references:ID"`
}
