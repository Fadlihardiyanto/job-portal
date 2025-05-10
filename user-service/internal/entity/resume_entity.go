package entity

type Resume struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	Attacment string `gorm:"column:attachment"`
	UserID    int    `gorm:"column:user_id"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
	Users     User   `gorm:"foreignKey:UserID;references:ID"`
}
