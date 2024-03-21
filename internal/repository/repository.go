package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

/*
Save
If you're updating many fields or have all the updated values available in a struct, Save may be faster because it generates a single SQL query to update all columns at once,
potentially reducing the overhead of multiple database operations.
*/
func (r *Repository[T]) Save(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

/*
Updates
If you're updating only a few specific fields, Updates may be faster because it generates a more targeted SQL query that affects fewer columns.
*/
func (r *Repository[T]) Updates(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Updates(entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}
