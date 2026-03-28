package service

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/repository"
)

var (
	ErrRegionIDsRequired        = errors.New("region_ids majburiy")
	ErrDistrictIDsRequired      = errors.New("district_ids majburiy")
	ErrDeliveryRegionNotFound   = errors.New("region id topilmadi")
	ErrDeliveryDistrictNotFound = errors.New("district id topilmadi")
	ErrDeliveryHierarchyInvalid = errors.New("tuman tanlangan viloyatlarga tegishli bo'lishi kerak")
)

type DeliveryAreasInput struct {
	RegionIDs   []uint `json:"region_ids"`
	DistrictIDs []uint `json:"district_ids"`
}

type DeliveryAreasResult struct {
	RegionIDs   []uint `json:"region_ids"`
	DistrictIDs []uint `json:"district_ids"`
}

type ContragentRegionDeliveryService interface {
	GetRegions() ([]adminDomain.Region, error)
	GetDistricts(regionID *uint) ([]adminDomain.District, error)
	GetMFYs(districtID *uint) ([]adminDomain.MFY, error)
	SaveDeliveryAreas(contragentID uint, input DeliveryAreasInput) (*DeliveryAreasResult, error)
	GetDeliveryAreas(contragentID uint) (*DeliveryAreasResult, error)
}

type contragentRegionDeliveryService struct {
	repo repository.ContragentRegionDeliveryRepository
}

func NewContragentRegionDeliveryService(repo repository.ContragentRegionDeliveryRepository) ContragentRegionDeliveryService {
	return &contragentRegionDeliveryService{repo: repo}
}

func (s *contragentRegionDeliveryService) GetRegions() ([]adminDomain.Region, error) {
	return s.repo.GetRegions()
}

func (s *contragentRegionDeliveryService) GetDistricts(regionID *uint) ([]adminDomain.District, error) {
	return s.repo.GetDistricts(regionID)
}

func (s *contragentRegionDeliveryService) GetMFYs(districtID *uint) ([]adminDomain.MFY, error) {
	return s.repo.GetMFYs(districtID)
}

func (s *contragentRegionDeliveryService) SaveDeliveryAreas(contragentID uint, input DeliveryAreasInput) (*DeliveryAreasResult, error) {
	regions := uniqueIDs(input.RegionIDs)
	districts := uniqueIDs(input.DistrictIDs)
	if len(regions) == 0 {
		return nil, ErrRegionIDsRequired
	}
	if len(districts) == 0 {
		return nil, ErrDistrictIDsRequired
	}

	regionSet := map[uint]struct{}{}
	for _, id := range regions {
		ok, err := s.repo.RegionExists(id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrDeliveryRegionNotFound
		}
		regionSet[id] = struct{}{}
	}
	for _, id := range districts {
		d, ok, err := s.repo.DistrictExists(id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrDeliveryDistrictNotFound
		}
		if _, exists := regionSet[d.RegionID]; !exists {
			return nil, ErrDeliveryHierarchyInvalid
		}
	}

	if err := s.repo.ReplaceDeliveryAreas(contragentID, regions, districts); err != nil {
		return nil, err
	}
	return &DeliveryAreasResult{
		RegionIDs:   regions,
		DistrictIDs: districts,
	}, nil
}

func (s *contragentRegionDeliveryService) GetDeliveryAreas(contragentID uint) (*DeliveryAreasResult, error) {
	regionIDs, districtIDs, err := s.repo.GetDeliveryAreas(contragentID)
	if err != nil {
		return nil, err
	}
	return &DeliveryAreasResult{
		RegionIDs:   regionIDs,
		DistrictIDs: districtIDs,
	}, nil
}

func uniqueIDs(ids []uint) []uint {
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
