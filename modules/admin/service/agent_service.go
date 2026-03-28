package service

import (
	"errors"
	"strings"

	"backend/internal/pkg/security"
	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var (
	ErrAgentNotFound       = errors.New("agent topilmadi")
	ErrAgentPhoneExists    = errors.New("telefon raqami allaqachon mavjud")
	ErrAgentHierarchy      = errors.New("viloyat, tuman va MFY mos kelmaydi")
	ErrAgentNameRequired   = errors.New("nom majburiy")
	ErrAgentLocationIDs    = errors.New("viloyat_id, tuman_id va mfy_id majburiy")
)

type AgentInput struct {
	Name                 string `json:"name"`
	ViloyatID            uint   `json:"viloyat_id"`
	TumanID              uint   `json:"tuman_id"`
	MFYID                uint   `json:"mfy_id"`
	Phone                string `json:"phone"`
	Status               string `json:"status"`
	PasswordSetupAllowed *bool  `json:"password_setup_allowed"`
	Password             string `json:"password"`
}

type AgentService interface {
	Create(input AgentInput) (*domain.Agent, error)
	GetPaginated(page, limit int) (*PaginatedAgents, error)
	GetByID(id uint) (*domain.Agent, error)
	Update(id uint, input AgentInput) (*domain.Agent, error)
	UpdateStatus(id uint, status string) (*domain.Agent, error)
	Delete(id uint) error
}

type PaginatedAgents struct {
	Items      []domain.Agent `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type agentService struct {
	repo       repository.AgentRepository
	regionRepo repository.RegionRepository
}

func NewAgentService(
	repo repository.AgentRepository,
	regionRepo repository.RegionRepository,
) AgentService {
	return &agentService{repo: repo, regionRepo: regionRepo}
}

func (s *agentService) Create(input AgentInput) (*domain.Agent, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrAgentNameRequired
	}
	if err := validateAgentPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err := s.validateLocation(input.ViloyatID, input.TumanID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrAgentPhoneExists
	}

	setupAllowed := true
	if input.PasswordSetupAllowed != nil {
		setupAllowed = *input.PasswordSetupAllowed
	}

	row := &domain.Agent{
		Name:                 input.Name,
		RegionID:             input.ViloyatID,
		DistrictID:           input.TumanID,
		MFYID:                input.MFYID,
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

func (s *agentService) GetPaginated(page, limit int) (*PaginatedAgents, error) {
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

	return &PaginatedAgents{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *agentService) GetByID(id uint) (*domain.Agent, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAgentNotFound
	}
	return row, nil
}

func (s *agentService) Update(id uint, input AgentInput) (*domain.Agent, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrAgentNameRequired
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAgentNotFound
	}

	if err = validateAgentPhone(input.Phone); err != nil {
		return nil, err
	}
	st, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err = s.validateLocation(input.ViloyatID, input.TumanID, input.MFYID); err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsByPhone(input.Phone, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrAgentPhoneExists
	}

	row.Name = input.Name
	row.RegionID = input.ViloyatID
	row.DistrictID = input.TumanID
	row.MFYID = input.MFYID
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

func (s *agentService) UpdateStatus(id uint, status string) (*domain.Agent, error) {
	st, err := normalizeStatus(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAgentNotFound
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *agentService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrAgentNotFound
	}
	return s.repo.Delete(id)
}

func validateAgentPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func (s *agentService) validateLocation(viloyatID, tumanID, mfyID uint) error {
	if viloyatID == 0 || tumanID == 0 || mfyID == 0 {
		return ErrAgentLocationIDs
	}

	district, err := s.regionRepo.GetDistrictByID(tumanID)
	if err != nil {
		return err
	}
	if district == nil {
		return ErrAgentHierarchy
	}
	if district.RegionID != viloyatID {
		return ErrAgentHierarchy
	}

	mfy, err := s.regionRepo.GetMFYByID(mfyID)
	if err != nil {
		return err
	}
	if mfy == nil {
		return ErrAgentHierarchy
	}
	if mfy.DistrictID != tumanID {
		return ErrAgentHierarchy
	}

	return nil
}

func (i *AgentInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Phone = strings.TrimSpace(i.Phone)
	i.Status = strings.TrimSpace(i.Status)
	i.Password = strings.TrimSpace(i.Password)
}
