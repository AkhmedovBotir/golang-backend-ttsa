package service

import (
	"errors"
	"regexp"
	"strings"

	"backend/internal/pkg/security"
	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var (
	ErrNeighborhoodShopNotFound       = errors.New("maxalla do'koni topilmadi")
	ErrNeighborhoodShopPhoneExists    = errors.New("telefon raqami allaqachon mavjud")
	ErrNeighborhoodShopInvalidINN       = errors.New("STIR (INN) 9 yoki 12 raqamdan iborat bo'lishi kerak yoki bo'sh qoldiring")
	ErrNeighborhoodShopHierarchy    = errors.New("viloyat, tuman va MFY mos kelmaydi")
	ErrNeighborhoodShopNameRequired = errors.New("nom majburiy")
	ErrNeighborhoodShopLocationIDs  = errors.New("region_id, district_id va mfy_id majburiy")
)

var neighborhoodShopInnRegex = regexp.MustCompile(`^(\d{9}|\d{12})$`)

type NeighborhoodShopInput struct {
	Name                 string `json:"name"`
	INN                  string `json:"inn"`
	RegionID             uint   `json:"region_id"`
	DistrictID           uint   `json:"district_id"`
	MFYID                uint   `json:"mfy_id"`
	Phone                string `json:"phone"`
	Logo                 string `json:"logo"`
	Status               string `json:"status"`
	PasswordSetupAllowed *bool  `json:"password_setup_allowed"`
	Password             string `json:"password"`
}

type NeighborhoodShopService interface {
	Create(input NeighborhoodShopInput) (*domain.NeighborhoodShop, error)
	GetPaginated(page, limit int) (*PaginatedNeighborhoodShops, error)
	GetByID(id uint) (*domain.NeighborhoodShop, error)
	Update(id uint, input NeighborhoodShopInput) (*domain.NeighborhoodShop, error)
	UpdateStatus(id uint, status string) (*domain.NeighborhoodShop, error)
	Delete(id uint) error
}

type PaginatedNeighborhoodShops struct {
	Items      []domain.NeighborhoodShop `json:"items"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

type neighborhoodShopService struct {
	repo       repository.NeighborhoodShopRepository
	regionRepo repository.RegionRepository
}

func NewNeighborhoodShopService(
	repo repository.NeighborhoodShopRepository,
	regionRepo repository.RegionRepository,
) NeighborhoodShopService {
	return &neighborhoodShopService{repo: repo, regionRepo: regionRepo}
}

func (s *neighborhoodShopService) Create(input NeighborhoodShopInput) (*domain.NeighborhoodShop, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrNeighborhoodShopNameRequired
	}
	if err := s.validateOptionalINN(input.INN); err != nil {
		return nil, err
	}
	if err := validateNeighborhoodShopPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err := s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrNeighborhoodShopPhoneExists
	}

	setupAllowed := true
	if input.PasswordSetupAllowed != nil {
		setupAllowed = *input.PasswordSetupAllowed
	}

	row := &domain.NeighborhoodShop{
		Name:                 input.Name,
		INN:                  input.INN,
		RegionID:             input.RegionID,
		DistrictID:           input.DistrictID,
		MFYID:                input.MFYID,
		Phone:                input.Phone,
		Logo:                 input.Logo,
		Status:               st,
		Password:             "",
		PasswordSetupAllowed: setupAllowed,
	}

	if err = s.repo.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *neighborhoodShopService) GetPaginated(page, limit int) (*PaginatedNeighborhoodShops, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := s.repo.GetPaginated(page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedNeighborhoodShops{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *neighborhoodShopService) GetByID(id uint) (*domain.NeighborhoodShop, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNeighborhoodShopNotFound
	}
	return row, nil
}

func (s *neighborhoodShopService) Update(id uint, input NeighborhoodShopInput) (*domain.NeighborhoodShop, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrNeighborhoodShopNameRequired
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNeighborhoodShopNotFound
	}

	if err = s.validateOptionalINN(input.INN); err != nil {
		return nil, err
	}
	if err = validateNeighborhoodShopPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err = s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrNeighborhoodShopPhoneExists
	}

	row.Name = input.Name
	row.INN = input.INN
	row.RegionID = input.RegionID
	row.DistrictID = input.DistrictID
	row.MFYID = input.MFYID
	row.Phone = input.Phone
	row.Logo = input.Logo
	row.Status = st

	if input.PasswordSetupAllowed != nil {
		row.PasswordSetupAllowed = *input.PasswordSetupAllowed
	}

	if strings.TrimSpace(input.Password) != "" {
		hashed, hashErr := security.HashPassword(strings.TrimSpace(input.Password))
		if hashErr != nil {
			return nil, hashErr
		}
		row.Password = hashed
	}

	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *neighborhoodShopService) UpdateStatus(id uint, status string) (*domain.NeighborhoodShop, error) {
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNeighborhoodShopNotFound
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *neighborhoodShopService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrNeighborhoodShopNotFound
	}
	return s.repo.Delete(id)
}

func (s *neighborhoodShopService) validateOptionalINN(inn string) error {
	inn = strings.TrimSpace(inn)
	if inn == "" {
		return nil
	}
	if !neighborhoodShopInnRegex.MatchString(inn) {
		return ErrNeighborhoodShopInvalidINN
	}
	return nil
}

func validateNeighborhoodShopPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func (s *neighborhoodShopService) validateLocation(regionID, districtID, mfyID uint) error {
	if regionID == 0 || districtID == 0 || mfyID == 0 {
		return ErrNeighborhoodShopLocationIDs
	}

	district, err := s.regionRepo.GetDistrictByID(districtID)
	if err != nil {
		return err
	}
	if district == nil {
		return ErrNeighborhoodShopHierarchy
	}
	if district.RegionID != regionID {
		return ErrNeighborhoodShopHierarchy
	}

	mfy, err := s.regionRepo.GetMFYByID(mfyID)
	if err != nil {
		return err
	}
	if mfy == nil {
		return ErrNeighborhoodShopHierarchy
	}
	if mfy.DistrictID != districtID {
		return ErrNeighborhoodShopHierarchy
	}

	return nil
}

func (i *NeighborhoodShopInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.INN = strings.TrimSpace(i.INN)
	i.Phone = strings.TrimSpace(i.Phone)
	i.Logo = strings.TrimSpace(i.Logo)
	i.Status = strings.TrimSpace(i.Status)
	i.Password = strings.TrimSpace(i.Password)
}
