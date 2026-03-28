package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"gorm.io/gorm"
)

type MarketplaceCategoryRepository interface {
	GetCategoryByID(id uint) (*adminDomain.Category, error)
	GetSubcategoryByID(id uint) (*adminDomain.Category, error)
	GetPaginatedCategories(page, limit int) ([]adminDomain.Category, int64, error)
	GetPaginatedSubcategories(page, limit int, parentID *uint) ([]adminDomain.Category, int64, error)
}

type marketplaceCategoryPostgresRepository struct {
	db *gorm.DB
}

func NewMarketplaceCategoryRepository(db *gorm.DB) MarketplaceCategoryRepository {
	return &marketplaceCategoryPostgresRepository{db: db}
}

func (r *marketplaceCategoryPostgresRepository) GetCategoryByID(id uint) (*adminDomain.Category, error) {
	var row adminDomain.Category
	err := r.db.
		Where("id = ? AND parent_id IS NULL AND status = ?", id, adminDomain.StatusActive).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceCategoryPostgresRepository) GetSubcategoryByID(id uint) (*adminDomain.Category, error) {
	var row adminDomain.Category
	err := r.db.
		Where("id = ? AND parent_id IS NOT NULL AND status = ?", id, adminDomain.StatusActive).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceCategoryPostgresRepository) GetPaginatedCategories(page, limit int) ([]adminDomain.Category, int64, error) {
	var rows []adminDomain.Category
	var total int64

	base := r.db.Model(&adminDomain.Category{}).
		Where("parent_id IS NULL AND status = ?", adminDomain.StatusActive)

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

func (r *marketplaceCategoryPostgresRepository) GetPaginatedSubcategories(page, limit int, parentID *uint) ([]adminDomain.Category, int64, error) {
	var rows []adminDomain.Category
	var total int64

	base := r.db.Model(&adminDomain.Category{}).
		Where("parent_id IS NOT NULL AND status = ?", adminDomain.StatusActive)
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
