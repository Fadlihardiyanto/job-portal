package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int          `gorm:"column:id;primaryKey"`
	FirstName       string       `gorm:"column:first_name"`
	LastName        string       `gorm:"column:last_name"`
	Email           string       `gorm:"column:email"`
	EmailToken      string       `gorm:"column:email_token"`
	EmailVerifiedAt sql.NullTime `gorm:"column:email_verified_at"`
	Password        string       `gorm:"column:password"`
	IsActive        bool         `gorm:"column:is_active"`
	Role            string       `gorm:"column:role"`
	About           string       `gorm:"column:about"`
	Photo           string       `gorm:"column:photo"`
	CreatedAt       time.Time    `gorm:"column:created_at"`
	UpdatedAt       time.Time    `gorm:"column:updated_at"`
	Token           string       `gorm:"-"`
}
