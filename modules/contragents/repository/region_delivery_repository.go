package repository

import (
	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type ContragentRegionDeliveryRepository interface {
	GetRegions() ([]adminDomain.Region, error)
	GetDistricts(regionID *uint) ([]adminDomain.District, error)
	GetMFYs(districtID *uint) ([]adminDomain.MFY, error)
	RegionExists(id uint) (bool, error)
	DistrictExists(id uint) (*adminDomain.District, bool, error)
	ReplaceDeliveryAreas(contragentID uint, regionIDs, districtIDs []uint) error
	GetDeliveryAreas(contragentID uint) ([]uint, []uint, error)
}

type contragentRegionDeliveryPostgresRepository struct {
	db *gorm.DB
}

func NewContragentRegionDeliveryRepository(db *gorm.DB) ContragentRegionDeliveryRepository {
	return &contragentRegionDeliveryPostgresRepository{db: db}
}

func (r *contragentRegionDeliveryPostgresRepository) GetRegions() ([]adminDomain.Region, error) {
	var rows []adminDomain.Region
	return rows, r.db.Order("id asc").Find(&rows).Error
}

func (r *contragentRegionDeliveryPostgresRepository) GetDistricts(regionID *uint) ([]adminDomain.District, error) {
	var rows []adminDomain.District
	q := r.db.Order("id asc")
	if regionID != nil {
		q = q.Where("region_id = ?", *regionID)
	}
	return rows, q.Find(&rows).Error
}

func (r *contragentRegionDeliveryPostgresRepository) GetMFYs(districtID *uint) ([]adminDomain.MFY, error) {
	var rows []adminDomain.MFY
	q := r.db.Order("id asc")
	if districtID != nil {
		q = q.Where("district_id = ?", *districtID)
	}
	return rows, q.Find(&rows).Error
}

func (r *contragentRegionDeliveryPostgresRepository) RegionExists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&adminDomain.Region{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *contragentRegionDeliveryPostgresRepository) DistrictExists(id uint) (*adminDomain.District, bool, error) {
	var row adminDomain.District
	err := r.db.First(&row, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &row, true, nil
}

func (r *contragentRegionDeliveryPostgresRepository) ReplaceDeliveryAreas(contragentID uint, regionIDs, districtIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("contragent_id = ?", contragentID).Delete(&domain.ContragentDeliveryRegion{}).Error; err != nil {
			return err
		}
		if err := tx.Where("contragent_id = ?", contragentID).Delete(&domain.ContragentDeliveryDistrict{}).Error; err != nil {
			return err
		}
		for _, id := range uniqueUint(regionIDs) {
			row := domain.ContragentDeliveryRegion{ContragentID: contragentID, RegionID: id}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		for _, id := range uniqueUint(districtIDs) {
			row := domain.ContragentDeliveryDistrict{ContragentID: contragentID, DistrictID: id}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *contragentRegionDeliveryPostgresRepository) GetDeliveryAreas(contragentID uint) ([]uint, []uint, error) {
	var regions []domain.ContragentDeliveryRegion
	var districts []domain.ContragentDeliveryDistrict
	if err := r.db.Where("contragent_id = ?", contragentID).Order("id asc").Find(&regions).Error; err != nil {
		return nil, nil, err
	}
	if err := r.db.Where("contragent_id = ?", contragentID).Order("id asc").Find(&districts).Error; err != nil {
		return nil, nil, err
	}
	regionIDs := make([]uint, 0, len(regions))
	for _, r := range regions {
		regionIDs = append(regionIDs, r.RegionID)
	}
	districtIDs := make([]uint, 0, len(districts))
	for _, d := range districts {
		districtIDs = append(districtIDs, d.DistrictID)
	}
	return regionIDs, districtIDs, nil
}

func uniqueUint(ids []uint) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
