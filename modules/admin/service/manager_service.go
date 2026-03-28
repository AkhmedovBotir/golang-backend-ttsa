package service

import (
	"errors"
	"strings"

	"backend/internal/pkg/security"
	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var (
	ErrManagerNotFound     = errors.New("menejer topilmadi")
	ErrManagerPhoneExists  = errors.New("telefon raqami allaqachon mavjud")
	ErrManagerHierarchy    = errors.New("viloyat topilmadi")
	ErrManagerNameRequired = errors.New("nom majburiy")
	ErrManagerRegionID     = errors.New("viloyat_id majburiy")
)

type ManagerInput struct {
	Name                 string `json:"name"`
	ViloyatID            uint   `json:"viloyat_id"`
	Phone                string `json:"phone"`
	Status               string `json:"status"`
	PasswordSetupAllowed *bool  `json:"password_setup_allowed"`
	Password             string `json:"password"`
}

type ManagerService interface {
	Create(input ManagerInput) (*domain.Manager, error)
	GetPaginated(page, limit int) (*PaginatedManagers, error)
	GetByID(id uint) (*domain.Manager, error)
	Update(id uint, input ManagerInput) (*domain.Manager, error)
	UpdateStatus(id uint, status string) (*domain.Manager, error)
	Delete(id uint) error
}

type PaginatedManagers struct {
	Items      []domain.Manager `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type managerService struct {
	repo       repository.ManagerRepository
	regionRepo repository.RegionRepository
}

func NewManagerService(
	repo repository.ManagerRepository,
	regionRepo repository.RegionRepository,
) ManagerService {
	return &managerService{repo: repo, regionRepo: regionRepo}
}

func (s *managerService) Create(input ManagerInput) (*domain.Manager, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrManagerNameRequired
	}
	if err := validateManagerPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err := s.validateViloyat(input.ViloyatID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrManagerPhoneExists
	}

	setupAllowed := true
	if input.PasswordSetupAllowed != nil {
		setupAllowed = *input.PasswordSetupAllowed
	}

	row := &domain.Manager{
		Name:                 input.Name,
		RegionID:             input.ViloyatID,
		Phone:                input.Phone,
		Status:               st,
		Password:             "",
		PasswordSetupAllowed: setupAllowed,
	}

	if err = s.repo.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *managerService) GetPaginated(page, limit int) (*PaginatedManagers, error) {
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

	return &PaginatedManagers{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *managerService) GetByID(id uint) (*domain.Manager, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrManagerNotFound
	}
	return row, nil
}

func (s *managerService) Update(id uint, input ManagerInput) (*domain.Manager, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrManagerNameRequired
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrManagerNotFound
	}

	if err = validateManagerPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err = s.validateViloyat(input.ViloyatID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrManagerPhoneExists
	}

	row.Name = input.Name
	row.RegionID = input.ViloyatID
	row.Phone = input.Phone
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

func (s *managerService) UpdateStatus(id uint, status string) (*domain.Manager, error) {
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrManagerNotFound
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *managerService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrManagerNotFound
	}
	return s.repo.Delete(id)
}

func validateManagerPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func (s *managerService) validateViloyat(viloyatID uint) error {
	if viloyatID == 0 {
		return ErrManagerRegionID
	}

	region, err := s.regionRepo.GetRegionByID(viloyatID)
	if err != nil {
		return err
	}
	if region == nil {
		return ErrManagerHierarchy
	}
	return nil
}

func (i *ManagerInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Phone = strings.TrimSpace(i.Phone)
	i.Status = strings.TrimSpace(i.Status)
	i.Password = strings.TrimSpace(i.Password)
}
