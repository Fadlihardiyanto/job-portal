package entity

import "time"

type Company struct {
	ID               int       `gorm:"column:id;primaryKey"`
	Name             string    `gorm:"column:name"`
	City             string    `gorm:"column:city"`
	OrganizationSize string    `gorm:"column:organization_size"`
	Logo             string    `gorm:"column:logo"`
	UserAccess       string    `gorm:"column:user_access"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
