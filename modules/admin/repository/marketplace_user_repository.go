package repository

import (
	"errors"
	"strings"

	"backend/modules/admin/domain"
	mpDomain "backend/modules/marketplace/domain"
	"gorm.io/gorm"
)

type MarketplaceUserFilter struct {
	Status     *string
	RegionID   *uint
	DistrictID *uint
	MFYID      *uint
	Phone      *string
	Query      *string
}

type MarketplaceUserRepository interface {
	GetPaginated(page, limit int, filter MarketplaceUserFilter) ([]mpDomain.User, int64, error)
	GetByID(id uint) (*mpDomain.User, error)
	Update(row *mpDomain.User) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
	GetRegionByID(id uint) (*domain.Region, error)
	GetDistrictByID(id uint) (*domain.District, error)
	GetMFYByID(id uint) (*domain.MFY, error)
}

type marketplaceUserPostgresRepository struct {
	db *gorm.DB
}

func NewMarketplaceUserRepository(db *gorm.DB) MarketplaceUserRepository {
	return &marketplaceUserPostgresRepository{db: db}
}

func (r *marketplaceUserPostgresRepository) GetPaginated(page, limit int, filter MarketplaceUserFilter) ([]mpDomain.User, int64, error) {
	var rows []mpDomain.User
	var total int64

	q := r.db.Model(&mpDomain.User{})
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.RegionID != nil {
		q = q.Where("region_id = ?", *filter.RegionID)
	}
	if filter.DistrictID != nil {
		q = q.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.MFYID != nil {
		q = q.Where("mfy_id = ?", *filter.MFYID)
	}
	if filter.Phone != nil && strings.TrimSpace(*filter.Phone) != "" {
		q = q.Where("phone = ?", strings.TrimSpace(*filter.Phone))
	}
	if filter.Query != nil && strings.TrimSpace(*filter.Query) != "" {
		raw := "%" + strings.TrimSpace(*filter.Query) + "%"
		q = q.Where("first_name ILIKE ? OR last_name ILIKE ?", raw, raw)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Order("id desc").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *marketplaceUserPostgresRepository) GetByID(id uint) (*mpDomain.User, error) {
	var row mpDomain.User
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceUserPostgresRepository) Update(row *mpDomain.User) error {
	return r.db.Save(row).Error
}

func (r *marketplaceUserPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&mpDomain.User{}, id).Error
}

func (r *marketplaceUserPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&mpDomain.User{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *marketplaceUserPostgresRepository) GetRegionByID(id uint) (*domain.Region, error) {
	var row domain.Region
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceUserPostgresRepository) GetDistrictByID(id uint) (*domain.District, error) {
	var row domain.District
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *marketplaceUserPostgresRepository) GetMFYByID(id uint) (*domain.MFY, error) {
	var row domain.MFY
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
