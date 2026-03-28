package service

import (
	"errors"
	"strings"

	"backend/internal/pkg/security"
	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var (
	ErrPunktNotFound     = errors.New("punkt topilmadi")
	ErrPunktPhoneExists  = errors.New("telefon raqami allaqachon mavjud")
	ErrPunktHierarchy    = errors.New("viloyat va tuman mos kelmaydi")
	ErrPunktNameRequired = errors.New("nom majburiy")
	ErrPunktLocationIDs  = errors.New("viloyat_id va tuman_id majburiy")
)

type PunktInput struct {
	Name                 string `json:"name"`
	ViloyatID            uint   `json:"viloyat_id"`
	TumanID              uint   `json:"tuman_id"`
	Phone                string `json:"phone"`
	Status               string `json:"status"`
	PasswordSetupAllowed *bool  `json:"password_setup_allowed"`
	Password             string `json:"password"`
}

type PunktService interface {
	Create(input PunktInput) (*domain.Punkt, error)
	GetPaginated(page, limit int) (*PaginatedPunkts, error)
	GetByID(id uint) (*domain.Punkt, error)
	Update(id uint, input PunktInput) (*domain.Punkt, error)
	UpdateStatus(id uint, status string) (*domain.Punkt, error)
	Delete(id uint) error
}

type PaginatedPunkts struct {
	Items      []domain.Punkt `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type punktService struct {
	repo       repository.PunktRepository
	regionRepo repository.RegionRepository
}

func NewPunktService(
	repo repository.PunktRepository,
	regionRepo repository.RegionRepository,
) PunktService {
	return &punktService{repo: repo, regionRepo: regionRepo}
}

func (s *punktService) Create(input PunktInput) (*domain.Punkt, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrPunktNameRequired
	}
	if err := validatePunktPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err := s.validateLocation(input.ViloyatID, input.TumanID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrPunktPhoneExists
	}

	setupAllowed := true
	if input.PasswordSetupAllowed != nil {
		setupAllowed = *input.PasswordSetupAllowed
	}

	row := &domain.Punkt{
		Name:                 input.Name,
		RegionID:             input.ViloyatID,
		DistrictID:           input.TumanID,
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

func (s *punktService) GetPaginated(page, limit int) (*PaginatedPunkts, error) {
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

	return &PaginatedPunkts{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *punktService) GetByID(id uint) (*domain.Punkt, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrPunktNotFound
	}
	return row, nil
}

func (s *punktService) Update(id uint, input PunktInput) (*domain.Punkt, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrPunktNameRequired
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrPunktNotFound
	}

	if err = validatePunktPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err = s.validateLocation(input.ViloyatID, input.TumanID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrPunktPhoneExists
	}

	row.Name = input.Name
	row.RegionID = input.ViloyatID
	row.DistrictID = input.TumanID
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

func (s *punktService) UpdateStatus(id uint, status string) (*domain.Punkt, error) {
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrPunktNotFound
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *punktService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrPunktNotFound
	}
	return s.repo.Delete(id)
}

func validatePunktPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func (s *punktService) validateLocation(viloyatID, tumanID uint) error {
	if viloyatID == 0 || tumanID == 0 {
		return ErrPunktLocationIDs
	}

	district, err := s.regionRepo.GetDistrictByID(tumanID)
	if err != nil {
		return err
	}
	if district == nil {
		return ErrPunktHierarchy
	}
	if district.RegionID != viloyatID {
		return ErrPunktHierarchy
	}

	return nil
}

func (i *PunktInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Phone = strings.TrimSpace(i.Phone)
	i.Status = strings.TrimSpace(i.Status)
	i.Password = strings.TrimSpace(i.Password)
}
