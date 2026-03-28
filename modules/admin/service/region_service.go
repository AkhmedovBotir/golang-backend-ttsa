package service

import (
	"errors"
	"strings"

	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var ErrRegionNotFound = errors.New("viloyat topilmadi")
var ErrDistrictNotFound = errors.New("tuman topilmadi")
var ErrMFYNotFound = errors.New("MFY topilmadi")
var ErrInvalidStatusField = errors.New("status noto'g'ri")
var ErrRegionIDRequired = errors.New("region_id majburiy")
var ErrDistrictIDRequired = errors.New("district_id majburiy")

type RegionService interface {
	CreateRegion(input RegionInput) (*domain.Region, error)
	GetRegions() ([]domain.Region, error)
	GetRegionByID(id uint) (*domain.Region, error)
	UpdateRegion(id uint, input RegionInput) (*domain.Region, error)
	UpdateRegionStatus(id uint, status string) (*domain.Region, error)
	DeleteRegion(id uint) error

	CreateDistrict(input DistrictInput) (*domain.District, error)
	GetDistricts() ([]domain.District, error)
	GetDistrictByID(id uint) (*domain.District, error)
	UpdateDistrict(id uint, input DistrictInput) (*domain.District, error)
	UpdateDistrictStatus(id uint, status string) (*domain.District, error)
	DeleteDistrict(id uint) error

	CreateMFY(input MFYInput) (*domain.MFY, error)
	GetMFYs() ([]domain.MFY, error)
	GetMFYByID(id uint) (*domain.MFY, error)
	UpdateMFY(id uint, input MFYInput) (*domain.MFY, error)
	UpdateMFYStatus(id uint, status string) (*domain.MFY, error)
	DeleteMFY(id uint) error
}

type regionService struct {
	repo repository.RegionRepository
}

func NewRegionService(repo repository.RegionRepository) RegionService {
	return &regionService{repo: repo}
}

type RegionInput struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

type DistrictInput struct {
	RegionID uint   `json:"region_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

type MFYInput struct {
	DistrictID uint   `json:"district_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Status     string `json:"status"`
}

func normalizeStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	if status != domain.StatusActive && status != domain.StatusInactive {
		return "", ErrInvalidStatusField
	}
	return status, nil
}

func (s *regionService) CreateRegion(input RegionInput) (*domain.Region, error) {
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row := &domain.Region{
		ExternalID: "",
		Name:       strings.TrimSpace(input.Name),
		Code:       strings.TrimSpace(input.Code),
		Status:     st,
	}
	return row, s.repo.CreateRegion(row)
}

func (s *regionService) GetRegions() ([]domain.Region, error) { return s.repo.GetRegions() }

func (s *regionService) GetRegionByID(id uint) (*domain.Region, error) {
	row, err := s.repo.GetRegionByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrRegionNotFound
	}
	return row, nil
}

func (s *regionService) UpdateRegion(id uint, input RegionInput) (*domain.Region, error) {
	row, err := s.repo.GetRegionByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrRegionNotFound
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row.Name = strings.TrimSpace(input.Name)
	row.Code = strings.TrimSpace(input.Code)
	row.Status = st
	return row, s.repo.UpdateRegion(row)
}

func (s *regionService) DeleteRegion(id uint) error {
	row, err := s.repo.GetRegionByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrRegionNotFound
	}
	return s.repo.DeleteRegion(id)
}

func (s *regionService) UpdateRegionStatus(id uint, status string) (*domain.Region, error) {
	row, err := s.repo.GetRegionByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrRegionNotFound
	}
	st, err := normalizeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.UpdateRegion(row); err != nil {
		return nil, err
	}
	if err = s.repo.UpdateDistrictStatusByRegionID(id, st); err != nil {
		return nil, err
	}
	if err = s.repo.UpdateMFYStatusByRegionID(id, st); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *regionService) CreateDistrict(input DistrictInput) (*domain.District, error) {
	if input.RegionID == 0 {
		return nil, ErrRegionIDRequired
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row := &domain.District{
		ExternalID: "",
		RegionID:   input.RegionID,
		Name:       strings.TrimSpace(input.Name),
		Code:       strings.TrimSpace(input.Code),
		Status:     st,
	}
	return row, s.repo.CreateDistrict(row)
}

func (s *regionService) GetDistricts() ([]domain.District, error) { return s.repo.GetDistricts() }

func (s *regionService) GetDistrictByID(id uint) (*domain.District, error) {
	row, err := s.repo.GetDistrictByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrDistrictNotFound
	}
	return row, nil
}

func (s *regionService) UpdateDistrict(id uint, input DistrictInput) (*domain.District, error) {
	row, err := s.repo.GetDistrictByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrDistrictNotFound
	}
	if input.RegionID == 0 {
		return nil, ErrRegionIDRequired
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row.RegionID = input.RegionID
	row.Name = strings.TrimSpace(input.Name)
	row.Code = strings.TrimSpace(input.Code)
	row.Status = st
	return row, s.repo.UpdateDistrict(row)
}

func (s *regionService) DeleteDistrict(id uint) error {
	row, err := s.repo.GetDistrictByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrDistrictNotFound
	}
	return s.repo.DeleteDistrict(id)
}

func (s *regionService) UpdateDistrictStatus(id uint, status string) (*domain.District, error) {
	row, err := s.repo.GetDistrictByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrDistrictNotFound
	}
	st, err := normalizeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.UpdateDistrict(row); err != nil {
		return nil, err
	}
	if err = s.repo.UpdateMFYStatusByDistrictID(id, st); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *regionService) CreateMFY(input MFYInput) (*domain.MFY, error) {
	if input.DistrictID == 0 {
		return nil, ErrDistrictIDRequired
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row := &domain.MFY{
		ExternalID: "",
		DistrictID: input.DistrictID,
		Name:       strings.TrimSpace(input.Name),
		Code:       strings.TrimSpace(input.Code),
		Status:     st,
	}
	return row, s.repo.CreateMFY(row)
}

func (s *regionService) GetMFYs() ([]domain.MFY, error) { return s.repo.GetMFYs() }

func (s *regionService) GetMFYByID(id uint) (*domain.MFY, error) {
	row, err := s.repo.GetMFYByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMFYNotFound
	}
	return row, nil
}

func (s *regionService) UpdateMFY(id uint, input MFYInput) (*domain.MFY, error) {
	row, err := s.repo.GetMFYByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMFYNotFound
	}
	if input.DistrictID == 0 {
		return nil, ErrDistrictIDRequired
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row.DistrictID = input.DistrictID
	row.Name = strings.TrimSpace(input.Name)
	row.Code = strings.TrimSpace(input.Code)
	row.Status = st
	return row, s.repo.UpdateMFY(row)
}

func (s *regionService) DeleteMFY(id uint) error {
	row, err := s.repo.GetMFYByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrMFYNotFound
	}
	return s.repo.DeleteMFY(id)
}

func (s *regionService) UpdateMFYStatus(id uint, status string) (*domain.MFY, error) {
	row, err := s.repo.GetMFYByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMFYNotFound
	}
	st, err := normalizeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	return row, s.repo.UpdateMFY(row)
}
