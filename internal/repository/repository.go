package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T) error {
	return db.Find(entities).Error
}

func (r *Repository[T]) FindByCondition(db *gorm.DB, entity *T, condition string, args ...any) error {
	return db.Where(condition, args...).Find(entity).Error
}

func (r *Repository[T]) FindByConditionWithPagination(db *gorm.DB, entities *[]T, condition string, page, pageSize int, args ...any) error {
	return db.Where(condition, args...).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(entities).Error
}

func (r *Repository[T]) FindByConditionWithPaginationAndOrder(db *gorm.DB, entities *[]T, condition string, page, pageSize int, order string, args ...any) error {
	return db.Where(condition, args...).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order(order).
		Find(entities).Error
}

func (r *Repository[T]) FindOne(db *gorm.DB, entity *T, condition string, args ...any) error {
	return db.Where(condition, args...).First(entity).Error
}
