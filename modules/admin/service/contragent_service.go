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
	ErrContragentNotFound              = errors.New("kontragent topilmadi")
	ErrContragentPhoneExists           = errors.New("telefon raqami allaqachon mavjud")
	ErrInvalidINN                      = errors.New("STIR (INN) 9 yoki 12 raqamdan iborat bo'lishi kerak")
	ErrContragentHierarchy             = errors.New("viloyat, tuman va MFY mos kelmaydi")
	ErrContragentNameRequired          = errors.New("nom majburiy")
	ErrContragentActivityTypeID        = errors.New("activity_type_id majburiy")
	ErrContragentLocationIDs           = errors.New("region_id, district_id va mfy_id majburiy")
	ErrContragentActivityTypeNotFound  = errors.New("faoliyat turi (kontragent turi) topilmadi")
)

var innRegex = regexp.MustCompile(`^(\d{9}|\d{12})$`)

type ContragentInput struct {
	Name                 string `json:"name"`
	INN                  string `json:"inn"`
	RegionID             uint   `json:"region_id"`
	DistrictID           uint   `json:"district_id"`
	MFYID                uint   `json:"mfy_id"`
	Phone                string `json:"phone"`
	Logo                 string `json:"logo"`
	ActivityTypeID       uint   `json:"activity_type_id"`
	Status               string `json:"status"`
	PasswordSetupAllowed *bool  `json:"password_setup_allowed"`
	Password             string `json:"password"`
}

type ContragentService interface {
	Create(input ContragentInput) (*domain.Contragent, error)
	GetPaginated(page, limit int) (*PaginatedContragents, error)
	GetByID(id uint) (*domain.Contragent, error)
	Update(id uint, input ContragentInput) (*domain.Contragent, error)
	UpdateStatus(id uint, status string) (*domain.Contragent, error)
	Delete(id uint) error
}

type PaginatedContragents struct {
	Items      []domain.Contragent `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type contragentService struct {
	repo       repository.ContragentRepository
	regionRepo repository.RegionRepository
	ctRepo     repository.ContragentTypeRepository
}

func NewContragentService(
	repo repository.ContragentRepository,
	regionRepo repository.RegionRepository,
	ctRepo repository.ContragentTypeRepository,
) ContragentService {
	return &contragentService{repo: repo, regionRepo: regionRepo, ctRepo: ctRepo}
}

func (s *contragentService) Create(input ContragentInput) (*domain.Contragent, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrContragentNameRequired
	}
	if err := s.validateINN(input.INN); err != nil {
		return nil, err
	}
	if err := validateContragentPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if input.ActivityTypeID == 0 {
		return nil, ErrContragentActivityTypeID
	}
	if err := s.mustExistContragentType(input.ActivityTypeID); err != nil {
		return nil, err
	}
	if err := s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrContragentPhoneExists
	}

	setupAllowed := true
	if input.PasswordSetupAllowed != nil {
		setupAllowed = *input.PasswordSetupAllowed
	}

	row := &domain.Contragent{
		Name:                 input.Name,
		INN:                  input.INN,
		RegionID:             input.RegionID,
		DistrictID:           input.DistrictID,
		MFYID:                input.MFYID,
		Phone:                input.Phone,
		Logo:                 input.Logo,
		ActivityTypeID:       input.ActivityTypeID,
		Status:               st,
		Password:             "",
		PasswordSetupAllowed: setupAllowed,
	}

	if err = s.repo.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *contragentService) GetPaginated(page, limit int) (*PaginatedContragents, error) {
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

	return &PaginatedContragents{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *contragentService) GetByID(id uint) (*domain.Contragent, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentNotFound
	}
	return row, nil
}

func (s *contragentService) Update(id uint, input ContragentInput) (*domain.Contragent, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrContragentNameRequired
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentNotFound
	}

	if err = s.validateINN(input.INN); err != nil {
		return nil, err
	}
	if err = validateContragentPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if input.ActivityTypeID == 0 {
		return nil, ErrContragentActivityTypeID
	}
	if err = s.mustExistContragentType(input.ActivityTypeID); err != nil {
		return nil, err
	}
	if err = s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrContragentPhoneExists
	}

	row.Name = input.Name
	row.INN = input.INN
	row.RegionID = input.RegionID
	row.DistrictID = input.DistrictID
	row.MFYID = input.MFYID
	row.Phone = input.Phone
	row.Logo = input.Logo
	row.ActivityTypeID = input.ActivityTypeID
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

func (s *contragentService) UpdateStatus(id uint, status string) (*domain.Contragent, error) {
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentNotFound
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *contragentService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrContragentNotFound
	}
	return s.repo.Delete(id)
}

func (s *contragentService) validateINN(inn string) error {
	if !innRegex.MatchString(inn) {
		return ErrInvalidINN
	}
	return nil
}

func validateContragentPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func (s *contragentService) validateLocation(regionID, districtID, mfyID uint) error {
	if regionID == 0 || districtID == 0 || mfyID == 0 {
		return ErrContragentLocationIDs
	}

	district, err := s.regionRepo.GetDistrictByID(districtID)
	if err != nil {
		return err
	}
	if district == nil {
		return ErrContragentHierarchy
	}
	if district.RegionID != regionID {
		return ErrContragentHierarchy
	}

	mfy, err := s.regionRepo.GetMFYByID(mfyID)
	if err != nil {
		return err
	}
	if mfy == nil {
		return ErrContragentHierarchy
	}
	if mfy.DistrictID != districtID {
		return ErrContragentHierarchy
	}

	return nil
}

func (i *ContragentInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.INN = strings.TrimSpace(i.INN)
	i.Phone = strings.TrimSpace(i.Phone)
	i.Logo = strings.TrimSpace(i.Logo)
	i.Status = strings.TrimSpace(i.Status)
	i.Password = strings.TrimSpace(i.Password)
}

func (s *contragentService) mustExistContragentType(id uint) error {
	row, err := s.ctRepo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrContragentActivityTypeNotFound
	}
	return nil
}
