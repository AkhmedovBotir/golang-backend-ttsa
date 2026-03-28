package repository

import (
	"strings"

	adminDomain "backend/modules/admin/domain"
	contrDomain "backend/modules/contragents/domain"
	"backend/modules/marketplace/domain"
	"gorm.io/gorm"
)

type MarketplaceProductFilter struct {
	CategoryID    *uint
	SubcategoryID *uint
	ContragentID  *uint
	Status        *string
	Query         *string
}

type MarketplaceProductRepository interface {
	ListApprovedProducts(filter MarketplaceProductFilter, page, limit int) ([]contrDomain.Product, int64, error)
	GetImages(productIDs []uint) (map[uint][]string, error)
	GetDeliveryAreasByContragentIDs(contragentIDs []uint) (map[uint]domain.DeliveryAreas, error)
	GetApprovedProductByID(id uint) (*contrDomain.Product, error)
}

type marketplaceProductPostgresRepository struct {
	db *gorm.DB
}

func NewMarketplaceProductRepository(db *gorm.DB) MarketplaceProductRepository {
	return &marketplaceProductPostgresRepository{db: db}
}

func (r *marketplaceProductPostgresRepository) ListApprovedProducts(filter MarketplaceProductFilter, page, limit int) ([]contrDomain.Product, int64, error) {
	base := r.db.Model(&contrDomain.Product{}).
		Where("moderation_status = ?", contrDomain.ProductModerationApproved)

	if filter.Status != nil {
		base = base.Where("status = ?", *filter.Status)
	} else {
		base = base.Where("status = ?", adminDomain.StatusActive)
	}

	if filter.CategoryID != nil {
		base = base.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.SubcategoryID != nil {
		base = base.Where("subcategory_id = ?", *filter.SubcategoryID)
	}
	if filter.ContragentID != nil {
		base = base.Where("contragent_id = ?", *filter.ContragentID)
	}
	if filter.Query != nil && strings.TrimSpace(*filter.Query) != "" {
		q := "%" + strings.TrimSpace(*filter.Query) + "%"
		base = base.Where("name ILIKE ? OR description ILIKE ?", q, q)
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
	var rows []contrDomain.Product
	if err := base.Order("id desc").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *marketplaceProductPostgresRepository) GetApprovedProductByID(id uint) (*contrDomain.Product, error) {
	var row contrDomain.Product
	err := r.db.
		Where("id = ? AND moderation_status = ? AND status = ?", id, contrDomain.ProductModerationApproved, adminDomain.StatusActive).
		First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceProductPostgresRepository) GetImages(productIDs []uint) (map[uint][]string, error) {
	out := make(map[uint][]string)
	if len(productIDs) == 0 {
		return out, nil
	}
	var rows []contrDomain.ProductImage
	if err := r.db.
		Where("product_id IN ?", productIDs).
		Order("sort_order asc, id asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ProductID] = append(out[row.ProductID], row.Image)
	}
	return out, nil
}

func (r *marketplaceProductPostgresRepository) GetDeliveryAreasByContragentIDs(contragentIDs []uint) (map[uint]domain.DeliveryAreas, error) {
	out := make(map[uint]domain.DeliveryAreas)
	if len(contragentIDs) == 0 {
		return out, nil
	}

	// Regions
	var regionRows []struct {
		ContragentID uint
		RegionID     uint
	}
	if err := r.db.
		Table("contragent_delivery_regions").
		Select("contragent_id, region_id").
		Where("contragent_id IN ?", contragentIDs).
		Order("id asc").
		Scan(&regionRows).Error; err != nil {
		return nil, err
	}
	for _, row := range regionRows {
		areas := out[row.ContragentID]
		areas.RegionIDs = append(areas.RegionIDs, row.RegionID)
		out[row.ContragentID] = areas
	}

	// Districts
	var districtRows []struct {
		ContragentID uint
		DistrictID   uint
	}
	if err := r.db.
		Table("contragent_delivery_districts").
		Select("contragent_id, district_id").
		Where("contragent_id IN ?", contragentIDs).
		Order("id asc").
		Scan(&districtRows).Error; err != nil {
		return nil, err
	}
	for _, row := range districtRows {
		areas := out[row.ContragentID]
		areas.DistrictIDs = append(areas.DistrictIDs, row.DistrictID)
		out[row.ContragentID] = areas
	}

	return out, nil
}
