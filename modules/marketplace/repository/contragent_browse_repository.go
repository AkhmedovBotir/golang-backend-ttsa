package repository

import (
	"strings"

	adminDomain "backend/modules/admin/domain"
	contrDomain "backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type CatSubPairForContragent struct {
	ContragentID  uint
	CategoryID    uint
	SubcategoryID uint
}

type MarketplaceContragentBrowseRepository interface {
	ListActiveContragents(query *string, page, limit int) ([]adminDomain.Contragent, int64, error)
	ListApprovedProductsForContragent(contragentID uint, limit int) ([]contrDomain.Product, error)
	GetDistinctCategoryPairsForContragents(contragentIDs []uint) ([]CatSubPairForContragent, error)
	GetCategoriesByIDs(ids []uint) (map[uint]adminDomain.Category, error)
}

type marketplaceContragentBrowsePostgresRepository struct {
	db *gorm.DB
}

func NewMarketplaceContragentBrowseRepository(db *gorm.DB) MarketplaceContragentBrowseRepository {
	return &marketplaceContragentBrowsePostgresRepository{db: db}
}

func (r *marketplaceContragentBrowsePostgresRepository) ListActiveContragents(query *string, page, limit int) ([]adminDomain.Contragent, int64, error) {
	base := r.db.Model(&adminDomain.Contragent{}).Where("status = ?", adminDomain.StatusActive)
	if query != nil && strings.TrimSpace(*query) != "" {
		p := "%" + strings.TrimSpace(*query) + "%"
		base = base.Where("name ILIKE ? OR inn ILIKE ? OR phone ILIKE ?", p, p, p)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	var rows []adminDomain.Contragent
	if err := base.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *marketplaceContragentBrowsePostgresRepository) ListApprovedProductsForContragent(contragentID uint, limit int) ([]contrDomain.Product, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	var rows []contrDomain.Product
	err := r.db.Model(&contrDomain.Product{}).
		Where("contragent_id = ? AND moderation_status = ? AND status = ?", contragentID, contrDomain.ProductModerationApproved, adminDomain.StatusActive).
		Order("id desc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *marketplaceContragentBrowsePostgresRepository) GetDistinctCategoryPairsForContragents(contragentIDs []uint) ([]CatSubPairForContragent, error) {
	if len(contragentIDs) == 0 {
		return nil, nil
	}
	var rows []CatSubPairForContragent
	err := r.db.Model(&contrDomain.Product{}).
		Select("DISTINCT contragent_id, category_id, subcategory_id").
		Where("contragent_id IN ? AND moderation_status = ? AND status = ?", contragentIDs, contrDomain.ProductModerationApproved, adminDomain.StatusActive).
		Scan(&rows).Error
	return rows, err
}

func (r *marketplaceContragentBrowsePostgresRepository) GetCategoriesByIDs(ids []uint) (map[uint]adminDomain.Category, error) {
	out := make(map[uint]adminDomain.Category)
	if len(ids) == 0 {
		return out, nil
	}
	var rows []adminDomain.Category
	if err := r.db.
		Where("id IN ? AND status = ?", ids, adminDomain.StatusActive).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		out[rows[i].ID] = rows[i]
	}
	return out, nil
}
