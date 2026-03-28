package service

import (
	"errors"
	"strings"
	"time"

	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/repository"
)

var (
	ErrDeliveryAreaNameRequired = errors.New("hudud nomi majburiy")
	ErrDeliveryAreaNotFound     = errors.New("yetkazib berish hududi topilmadi")
	ErrDeliveryLocationIDs      = errors.New("viloyat_id, tuman_id va mfy_id majburiy")
	ErrDeliveryHierarchy        = errors.New("viloyat, tuman va mfy mos emas")
)

type DeliveryAreaInput struct {
	Name       string `json:"name"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	IsDefault  bool   `json:"is_default"`
}

type DeliveryAreaOutput struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	RegionID   uint      `json:"region_id"`
	DistrictID uint      `json:"district_id"`
	MFYID      uint      `json:"mfy_id"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DeliveryAreaService interface {
	Create(userID uint, input DeliveryAreaInput) (*DeliveryAreaOutput, error)
	List(userID uint) ([]DeliveryAreaOutput, error)
	Update(userID, id uint, input DeliveryAreaInput) (*DeliveryAreaOutput, error)
	Delete(userID, id uint) error
	SetDefault(userID, id uint) (*DeliveryAreaOutput, error)
}

type deliveryAreaService struct {
	repo repository.DeliveryAreaRepository
}

func NewDeliveryAreaService(repo repository.DeliveryAreaRepository) DeliveryAreaService {
	return &deliveryAreaService{repo: repo}
}

func (s *deliveryAreaService) Create(userID uint, input DeliveryAreaInput) (*DeliveryAreaOutput, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := s.validateInput(input); err != nil {
		return nil, err
	}
	row := &domain.DeliveryArea{
		UserID:     userID,
		Name:       input.Name,
		RegionID:   input.RegionID,
		DistrictID: input.DistrictID,
		MFYID:      input.MFYID,
		IsDefault:  input.IsDefault,
	}
	count, err := s.repo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		row.IsDefault = true
	}
	if row.IsDefault {
		if err = s.repo.ResetDefaultByUserID(userID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Create(row); err != nil {
		return nil, err
	}
	out := mapDeliveryArea(row)
	return &out, nil
}

func (s *deliveryAreaService) List(userID uint) ([]DeliveryAreaOutput, error) {
	rows, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]DeliveryAreaOutput, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapDeliveryArea(&row))
	}
	return out, nil
}

func (s *deliveryAreaService) Update(userID, id uint, input DeliveryAreaInput) (*DeliveryAreaOutput, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := s.validateInput(input); err != nil {
		return nil, err
	}
	row, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrDeliveryAreaNotFound
	}
	row.Name = input.Name
	row.RegionID = input.RegionID
	row.DistrictID = input.DistrictID
	row.MFYID = input.MFYID
	row.IsDefault = input.IsDefault
	if row.IsDefault {
		if err = s.repo.ResetDefaultByUserID(userID); err != nil {
			return nil, err
		}
	}
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	out := mapDeliveryArea(row)
	return &out, nil
}

func (s *deliveryAreaService) Delete(userID, id uint) error {
	row, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrDeliveryAreaNotFound
	}
	return s.repo.Delete(row)
}

func (s *deliveryAreaService) SetDefault(userID, id uint) (*DeliveryAreaOutput, error) {
	row, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrDeliveryAreaNotFound
	}
	if err = s.repo.ResetDefaultByUserID(userID); err != nil {
		return nil, err
	}
	row.IsDefault = true
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	out := mapDeliveryArea(row)
	return &out, nil
}

func (s *deliveryAreaService) validateInput(input DeliveryAreaInput) error {
	if input.Name == "" {
		return ErrDeliveryAreaNameRequired
	}
	if input.RegionID == 0 || input.DistrictID == 0 || input.MFYID == 0 {
		return ErrDeliveryLocationIDs
	}
	district, err := s.repo.GetDistrictByID(input.DistrictID)
	if err != nil {
		return err
	}
	if district == nil || district.RegionID != input.RegionID {
		return ErrDeliveryHierarchy
	}
	mfy, err := s.repo.GetMFYByID(input.MFYID)
	if err != nil {
		return err
	}
	if mfy == nil || mfy.DistrictID != input.DistrictID {
		return ErrDeliveryHierarchy
	}
	region, err := s.repo.GetRegionByID(input.RegionID)
	if err != nil {
		return err
	}
	if region == nil {
		return ErrDeliveryHierarchy
	}
	return nil
}

func mapDeliveryArea(row *domain.DeliveryArea) DeliveryAreaOutput {
	return DeliveryAreaOutput{
		ID:         row.ID,
		Name:       row.Name,
		RegionID:   row.RegionID,
		DistrictID: row.DistrictID,
		MFYID:      row.MFYID,
		IsDefault:  row.IsDefault,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
