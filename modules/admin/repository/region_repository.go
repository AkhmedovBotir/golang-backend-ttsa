package repository

import (
	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type RegionRepository interface {
	CreateRegion(region *domain.Region) error
	GetRegions() ([]domain.Region, error)
	GetRegionByID(id uint) (*domain.Region, error)
	UpdateRegion(region *domain.Region) error
	DeleteRegion(id uint) error

	CreateDistrict(district *domain.District) error
	GetDistricts() ([]domain.District, error)
	GetDistrictByID(id uint) (*domain.District, error)
	UpdateDistrict(district *domain.District) error
	UpdateDistrictStatusByRegionID(regionID uint, status string) error
	DeleteDistrict(id uint) error

	CreateMFY(mfy *domain.MFY) error
	GetMFYs() ([]domain.MFY, error)
	GetMFYByID(id uint) (*domain.MFY, error)
	UpdateMFY(mfy *domain.MFY) error
	UpdateMFYStatusByDistrictID(districtID uint, status string) error
	UpdateMFYStatusByRegionID(regionID uint, status string) error
	DeleteMFY(id uint) error
}

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{db: db}
}

func (r *regionRepository) CreateRegion(region *domain.Region) error { return r.db.Create(region).Error }
func (r *regionRepository) GetRegions() ([]domain.Region, error) {
	var rows []domain.Region
	return rows, r.db.Order("id asc").Find(&rows).Error
}
func (r *regionRepository) GetRegionByID(id uint) (*domain.Region, error) {
	var row domain.Region
	if err := r.db.First(&row, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
func (r *regionRepository) UpdateRegion(region *domain.Region) error { return r.db.Save(region).Error }
func (r *regionRepository) DeleteRegion(id uint) error               { return r.db.Delete(&domain.Region{}, id).Error }

func (r *regionRepository) CreateDistrict(district *domain.District) error { return r.db.Create(district).Error }
func (r *regionRepository) GetDistricts() ([]domain.District, error) {
	var rows []domain.District
	return rows, r.db.Order("id asc").Find(&rows).Error
}
func (r *regionRepository) GetDistrictByID(id uint) (*domain.District, error) {
	var row domain.District
	if err := r.db.First(&row, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
func (r *regionRepository) UpdateDistrict(district *domain.District) error { return r.db.Save(district).Error }
func (r *regionRepository) UpdateDistrictStatusByRegionID(regionID uint, status string) error {
	return r.db.Model(&domain.District{}).Where("region_id = ?", regionID).Update("status", status).Error
}
func (r *regionRepository) DeleteDistrict(id uint) error                   { return r.db.Delete(&domain.District{}, id).Error }

func (r *regionRepository) CreateMFY(mfy *domain.MFY) error { return r.db.Create(mfy).Error }
func (r *regionRepository) GetMFYs() ([]domain.MFY, error) {
	var rows []domain.MFY
	return rows, r.db.Order("id asc").Find(&rows).Error
}
func (r *regionRepository) GetMFYByID(id uint) (*domain.MFY, error) {
	var row domain.MFY
	if err := r.db.First(&row, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
func (r *regionRepository) UpdateMFY(mfy *domain.MFY) error { return r.db.Save(mfy).Error }
func (r *regionRepository) UpdateMFYStatusByDistrictID(districtID uint, status string) error {
	return r.db.Model(&domain.MFY{}).Where("district_id = ?", districtID).Update("status", status).Error
}
func (r *regionRepository) UpdateMFYStatusByRegionID(regionID uint, status string) error {
	subQuery := r.db.Model(&domain.District{}).Select("id").Where("region_id = ?", regionID)
	return r.db.Model(&domain.MFY{}).Where("district_id IN (?)", subQuery).Update("status", status).Error
}
func (r *regionRepository) DeleteMFY(id uint) error          { return r.db.Delete(&domain.MFY{}, id).Error }
