package repository

import (
	"errors"

	"backend/modules/admin/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(row *domain.Category) error
	GetByID(id uint) (*domain.Category, error)
	Update(row *domain.Category) error
	UpdateChildrenStatusByParentID(parentID uint, status string) error
	Delete(id uint) error
	ExistsBySlug(slug string, exceptID uint) (bool, error)
	ParentExists(id uint) (bool, error)
	GetPaginatedCategories(page, limit int) ([]domain.Category, int64, error)
	GetPaginatedSubcategories(page, limit int, parentID *uint) ([]domain.Category, int64, error)
}

type categoryPostgresRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryPostgresRepository{db: db}
}

func (r *categoryPostgresRepository) Create(row *domain.Category) error {
	return r.db.Create(row).Error
}

func (r *categoryPostgresRepository) GetByID(id uint) (*domain.Category, error) {
	var row domain.Category
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *categoryPostgresRepository) Update(row *domain.Category) error {
	return r.db.Save(row).Error
}

func (r *categoryPostgresRepository) UpdateChildrenStatusByParentID(parentID uint, status string) error {
	return r.db.Model(&domain.Category{}).
		Where("parent_id = ?", parentID).
		Update("status", status).Error
}

func (r *categoryPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *categoryPostgresRepository) ExistsBySlug(slug string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Category{}).Where("slug = ?", slug)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryPostgresRepository) ParentExists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Category{}).Where("id = ? AND parent_id IS NULL", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryPostgresRepository) GetPaginatedCategories(page, limit int) ([]domain.Category, int64, error) {
	var rows []domain.Category
	var total int64

	base := r.db.Model(&domain.Category{}).Where("parent_id IS NULL")
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := base.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *categoryPostgresRepository) GetPaginatedSubcategories(page, limit int, parentID *uint) ([]domain.Category, int64, error) {
	var rows []domain.Category
	var total int64

	base := r.db.Model(&domain.Category{}).Where("parent_id IS NOT NULL")
	if parentID != nil {
		base = base.Where("parent_id = ?", *parentID)
	}
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := base.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
