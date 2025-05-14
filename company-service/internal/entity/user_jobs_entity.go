package entity

import "time"

type UserJobs struct {
	ID        int       `gorm:"column:id;primaryKey"`
	UserID    int       `gorm:"column:user_id"`
	ResumeID  int       `gorm:"column:resume_id"`
	JobsID    int       `gorm:"column:jobs_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Users     User      `gorm:"foreignKey:UserID;references:ID"`
	Jobs      Jobs      `gorm:"foreignKey:JobsID;references:ID"`
	Resume    Resume    `gorm:"foreignKey:ResumeID;references:ID"`
}

func (UserJobs) TableName() string {
	return "user_jobs"
}
