package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"backend/modules/admin/repository"
	mpDomain "backend/modules/marketplace/domain"
)

var (
	ErrMarketplaceUserNotFound     = errors.New("marketplace user topilmadi")
	ErrMarketplaceUserPhoneExists  = errors.New("telefon raqami allaqachon mavjud")
	ErrMarketplaceUserNameRequired = errors.New("ism va familiya majburiy")
	ErrMarketplaceUserGender       = errors.New("jins noto'g'ri")
	ErrMarketplaceUserBirthDate    = errors.New("tug'ilgan sana noto'g'ri")
	ErrMarketplaceUserLocation     = errors.New("viloyat, tuman va mfy mos emas")
)

var marketplacePhoneRegex = regexp.MustCompile(`^\+998\d{9}$`)

type MarketplaceUserInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	BirthDate  string `json:"birth_date"` // YYYY-MM-DD
	Status     string `json:"status"`
}

type PaginatedMarketplaceUsers struct {
	Items      []mpDomain.User `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

type MarketplaceUserService interface {
	GetPaginated(page, limit int, filter repository.MarketplaceUserFilter) (*PaginatedMarketplaceUsers, error)
	GetByID(id uint) (*mpDomain.User, error)
	Update(id uint, input MarketplaceUserInput) (*mpDomain.User, error)
	UpdateStatus(id uint, status string) (*mpDomain.User, error)
	Delete(id uint) error
}

type marketplaceUserService struct {
	repo repository.MarketplaceUserRepository
}

func NewMarketplaceUserService(repo repository.MarketplaceUserRepository) MarketplaceUserService {
	return &marketplaceUserService{repo: repo}
}

func (s *marketplaceUserService) GetPaginated(page, limit int, filter repository.MarketplaceUserFilter) (*PaginatedMarketplaceUsers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	items, total, err := s.repo.GetPaginated(page, limit, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &PaginatedMarketplaceUsers{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *marketplaceUserService) GetByID(id uint) (*mpDomain.User, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMarketplaceUserNotFound
	}
	return row, nil
}

func (s *marketplaceUserService) Update(id uint, input MarketplaceUserInput) (*mpDomain.User, error) {
	input.normalize()
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMarketplaceUserNotFound
	}

	if input.FirstName == "" || input.LastName == "" {
		return nil, ErrMarketplaceUserNameRequired
	}
	if !marketplacePhoneRegex.MatchString(input.Phone) {
		return nil, ErrInvalidPhone
	}
	if input.Gender != mpDomain.GenderMale && input.Gender != mpDomain.GenderFemale {
		return nil, ErrMarketplaceUserGender
	}
	if input.RegionID == 0 || input.DistrictID == 0 || input.MFYID == 0 {
		return nil, ErrMarketplaceUserLocation
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, ErrMarketplaceUserBirthDate
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrMarketplaceUserPhoneExists
	}
	if err = s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}

	row.FirstName = input.FirstName
	row.LastName = input.LastName
	row.Gender = input.Gender
	row.Phone = input.Phone
	row.RegionID = input.RegionID
	row.DistrictID = input.DistrictID
	row.MFYID = input.MFYID
	row.BirthDate = birthDate
	row.Status = st

	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *marketplaceUserService) UpdateStatus(id uint, status string) (*mpDomain.User, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMarketplaceUserNotFound
	}
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *marketplaceUserService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrMarketplaceUserNotFound
	}
	return s.repo.Delete(id)
}

func (s *marketplaceUserService) validateLocation(regionID, districtID, mfyID uint) error {
	region, err := s.repo.GetRegionByID(regionID)
	if err != nil {
		return err
	}
	if region == nil {
		return ErrMarketplaceUserLocation
	}
	district, err := s.repo.GetDistrictByID(districtID)
	if err != nil {
		return err
	}
	if district == nil || district.RegionID != regionID {
		return ErrMarketplaceUserLocation
	}
	mfy, err := s.repo.GetMFYByID(mfyID)
	if err != nil {
		return err
	}
	if mfy == nil || mfy.DistrictID != districtID {
		return ErrMarketplaceUserLocation
	}
	return nil
}

func (i *MarketplaceUserInput) normalize() {
	i.FirstName = strings.TrimSpace(i.FirstName)
	i.LastName = strings.TrimSpace(i.LastName)
	i.Gender = strings.TrimSpace(i.Gender)
	i.Phone = strings.TrimSpace(i.Phone)
	i.BirthDate = strings.TrimSpace(i.BirthDate)
	i.Status = strings.TrimSpace(i.Status)
}
