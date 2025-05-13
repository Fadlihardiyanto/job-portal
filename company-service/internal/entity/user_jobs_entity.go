package entity

type UserJobs struct {
	ID        int    `gorm:"column:id;primaryKey"`
	UserID    int    `gorm:"column:user_id"`
	ResumeID  int    `gorm:"column:resume_id"`
	JobID     int    `gorm:"column:job_id"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
	Users     User   `gorm:"foreignKey:UserID;references:ID"`
	Jobs      Jobs   `gorm:"foreignKey:JobID;references:ID"`
	Resume    Resume `gorm:"foreignKey:ResumeID;references:ID"`
}
