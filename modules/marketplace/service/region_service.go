package service

import (
	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/repository"
)

type RegionService interface {
	GetRegions() ([]adminDomain.Region, error)
	GetDistricts(regionID *uint) ([]adminDomain.District, error)
	GetMFYs(districtID *uint) ([]adminDomain.MFY, error)
}

type regionService struct {
	repo repository.RegionRepository
}

func NewRegionService(repo repository.RegionRepository) RegionService {
	return &regionService{repo: repo}
}

func (s *regionService) GetRegions() ([]adminDomain.Region, error) {
	return s.repo.GetRegions()
}

func (s *regionService) GetDistricts(regionID *uint) ([]adminDomain.District, error) {
	return s.repo.GetDistricts(regionID)
}

func (s *regionService) GetMFYs(districtID *uint) ([]adminDomain.MFY, error) {
	return s.repo.GetMFYs(districtID)
}
