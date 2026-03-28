package service

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/repository"
)

var (
	ErrContragentCategoryNotFound    = errors.New("kategoriya topilmadi")
	ErrContragentSubcategoryNotFound = errors.New("subkategoriya topilmadi")
)

type ContragentPaginatedCategories struct {
	Items      []adminDomain.Category `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

type ContragentCategoryService interface {
	GetCategories(page, limit int) (*ContragentPaginatedCategories, error)
	GetCategoryByID(id uint) (*adminDomain.Category, error)
	GetSubcategories(page, limit int, parentID *uint) (*ContragentPaginatedCategories, error)
	GetSubcategoryByID(id uint) (*adminDomain.Category, error)
}

type contragentCategoryService struct {
	repo repository.ContragentCategoryRepository
}

func NewContragentCategoryService(repo repository.ContragentCategoryRepository) ContragentCategoryService {
	return &contragentCategoryService{repo: repo}
}

func (s *contragentCategoryService) GetCategories(page, limit int) (*ContragentPaginatedCategories, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	items, total, err := s.repo.GetPaginatedCategories(page, limit)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &ContragentPaginatedCategories{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *contragentCategoryService) GetCategoryByID(id uint) (*adminDomain.Category, error) {
	row, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentCategoryNotFound
	}
	return row, nil
}

func (s *contragentCategoryService) GetSubcategories(page, limit int, parentID *uint) (*ContragentPaginatedCategories, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	items, total, err := s.repo.GetPaginatedSubcategories(page, limit, parentID)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &ContragentPaginatedCategories{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *contragentCategoryService) GetSubcategoryByID(id uint) (*adminDomain.Category, error) {
	row, err := s.repo.GetSubcategoryByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentSubcategoryNotFound
	}
	return row, nil
}
