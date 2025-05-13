package entity

import "time"

type Resume struct {
	ID         int       `gorm:"column:id;primaryKey"`
	Name       string    `gorm:"column:name"`
	Attachment string    `gorm:"column:attachment"`
	UserID     int       `gorm:"column:user_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	Users      User      `gorm:"foreignKey:UserID;references:ID"`
}

func (Resume) TableName() string {
	return "resume"
}
