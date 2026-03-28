package repository

import (
	"strings"

	adminDomain "backend/modules/admin/domain"
	contrDomain "backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type MarketplaceSearchRepository interface {
	SearchProducts(q string, limit int) ([]contrDomain.Product, error)
	SearchMainCategories(q string, limit int) ([]adminDomain.Category, error)
	SearchSubcategories(q string, limit int) ([]adminDomain.Category, error)
	SearchContragents(q string, limit int) ([]adminDomain.Contragent, error)
}

type marketplaceSearchPostgresRepository struct {
	db *gorm.DB
}

func NewMarketplaceSearchRepository(db *gorm.DB) MarketplaceSearchRepository {
	return &marketplaceSearchPostgresRepository{db: db}
}

func clampSearchLimit(limit int) int {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return limit
}

func (r *marketplaceSearchPostgresRepository) SearchProducts(q string, limit int) ([]contrDomain.Product, error) {
	limit = clampSearchLimit(limit)
	pattern := "%" + strings.TrimSpace(q) + "%"
	var rows []contrDomain.Product
	err := r.db.Model(&contrDomain.Product{}).
		Where("moderation_status = ? AND status = ?", contrDomain.ProductModerationApproved, adminDomain.StatusActive).
		Where("name ILIKE ? OR description ILIKE ?", pattern, pattern).
		Order("id desc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *marketplaceSearchPostgresRepository) SearchMainCategories(q string, limit int) ([]adminDomain.Category, error) {
	limit = clampSearchLimit(limit)
	pattern := "%" + strings.TrimSpace(q) + "%"
	var rows []adminDomain.Category
	err := r.db.Model(&adminDomain.Category{}).
		Where("parent_id IS NULL AND status = ?", adminDomain.StatusActive).
		Where("name ILIKE ? OR slug ILIKE ?", pattern, pattern).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *marketplaceSearchPostgresRepository) SearchSubcategories(q string, limit int) ([]adminDomain.Category, error) {
	limit = clampSearchLimit(limit)
	pattern := "%" + strings.TrimSpace(q) + "%"
	var rows []adminDomain.Category
	err := r.db.Model(&adminDomain.Category{}).
		Where("parent_id IS NOT NULL AND status = ?", adminDomain.StatusActive).
		Where("name ILIKE ? OR slug ILIKE ?", pattern, pattern).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *marketplaceSearchPostgresRepository) SearchContragents(q string, limit int) ([]adminDomain.Contragent, error) {
	limit = clampSearchLimit(limit)
	pattern := "%" + strings.TrimSpace(q) + "%"
	var rows []adminDomain.Contragent
	err := r.db.Model(&adminDomain.Contragent{}).
		Where("status = ?", adminDomain.StatusActive).
		Where("name ILIKE ? OR inn ILIKE ? OR phone ILIKE ?", pattern, pattern, pattern).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}
