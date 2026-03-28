package service

import (
	"errors"
	"strings"

	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var ErrContragentTypeNotFound = errors.New("kontragent turi topilmadi")

type ContragentTypeInput struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Status string `json:"status"`
}

type ContragentTypeService interface {
	Create(input ContragentTypeInput) (*domain.ContragentType, error)
	GetAll() ([]domain.ContragentType, error)
	GetPaginated(page, limit int) (*PaginatedContragentTypes, error)
	GetByID(id uint) (*domain.ContragentType, error)
	Update(id uint, input ContragentTypeInput) (*domain.ContragentType, error)
	UpdateStatus(id uint, status string) (*domain.ContragentType, error)
	Delete(id uint) error
}

type PaginatedContragentTypes struct {
	Items      []domain.ContragentType `json:"items"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

type contragentTypeService struct {
	repo repository.ContragentTypeRepository
}

func NewContragentTypeService(repo repository.ContragentTypeRepository) ContragentTypeService {
	return &contragentTypeService{repo: repo}
}

func (s *contragentTypeService) Create(input ContragentTypeInput) (*domain.ContragentType, error) {
	status, err := normalizeContragentTypeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row := &domain.ContragentType{
		Name:   strings.TrimSpace(input.Name),
		Icon:   strings.TrimSpace(input.Icon),
		Status: status,
	}
	return row, s.repo.Create(row)
}

func (s *contragentTypeService) GetAll() ([]domain.ContragentType, error) {
	return s.repo.GetAll()
}

func (s *contragentTypeService) GetPaginated(page, limit int) (*PaginatedContragentTypes, error) {
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

	return &PaginatedContragentTypes{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *contragentTypeService) GetByID(id uint) (*domain.ContragentType, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentTypeNotFound
	}
	return row, nil
}

func (s *contragentTypeService) Update(id uint, input ContragentTypeInput) (*domain.ContragentType, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentTypeNotFound
	}

	status, err := normalizeContragentTypeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	row.Name = strings.TrimSpace(input.Name)
	row.Icon = strings.TrimSpace(input.Icon)
	row.Status = status
	return row, s.repo.Update(row)
}

func (s *contragentTypeService) UpdateStatus(id uint, status string) (*domain.ContragentType, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentTypeNotFound
	}
	normalized, err := normalizeContragentTypeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = normalized
	return row, s.repo.Update(row)
}

func (s *contragentTypeService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrContragentTypeNotFound
	}
	return s.repo.Delete(id)
}

func normalizeContragentTypeStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	if status != domain.StatusActive && status != domain.StatusInactive {
		return "", ErrInvalidStatusField
	}
	return status, nil
}
