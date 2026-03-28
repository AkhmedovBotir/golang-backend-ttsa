package repository

import (
	adminDomain "backend/modules/admin/domain"
	"gorm.io/gorm"
)

type RegionRepository interface {
	GetRegions() ([]adminDomain.Region, error)
	GetDistricts(regionID *uint) ([]adminDomain.District, error)
	GetMFYs(districtID *uint) ([]adminDomain.MFY, error)
}

type regionPostgresRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionPostgresRepository{db: db}
}

func (r *regionPostgresRepository) GetRegions() ([]adminDomain.Region, error) {
	var rows []adminDomain.Region
	return rows, r.db.Order("id asc").Find(&rows).Error
}

func (r *regionPostgresRepository) GetDistricts(regionID *uint) ([]adminDomain.District, error) {
	var rows []adminDomain.District
	q := r.db.Order("id asc")
	if regionID != nil {
		q = q.Where("region_id = ?", *regionID)
	}
	return rows, q.Find(&rows).Error
}

func (r *regionPostgresRepository) GetMFYs(districtID *uint) ([]adminDomain.MFY, error) {
	var rows []adminDomain.MFY
	q := r.db.Order("id asc")
	if districtID != nil {
		q = q.Where("district_id = ?", *districtID)
	}
	return rows, q.Find(&rows).Error
}
