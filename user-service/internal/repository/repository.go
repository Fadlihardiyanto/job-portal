package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) FindByID(db *gorm.DB, id string) (*T, error) {
	var entity T
	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) FindAll(db *gorm.DB) ([]T, error) {
	var entities []T
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}
