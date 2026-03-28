package service

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/repository"
)

var (
	ErrMarketplaceCategoryNotFound    = errors.New("kategoriya topilmadi")
	ErrMarketplaceSubcategoryNotFound = errors.New("subkategoriya topilmadi")
)

type MarketplacePaginatedCategories struct {
	Items      []adminDomain.Category `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

type MarketplaceCategoryService interface {
	GetCategories(page, limit int) (*MarketplacePaginatedCategories, error)
	GetCategoryByID(id uint) (*adminDomain.Category, error)
	GetSubcategories(page, limit int, parentID *uint) (*MarketplacePaginatedCategories, error)
	GetSubcategoryByID(id uint) (*adminDomain.Category, error)
}

type marketplaceCategoryService struct {
	repo repository.MarketplaceCategoryRepository
}

func NewMarketplaceCategoryService(repo repository.MarketplaceCategoryRepository) MarketplaceCategoryService {
	return &marketplaceCategoryService{repo: repo}
}

func (s *marketplaceCategoryService) GetCategories(page, limit int) (*MarketplacePaginatedCategories, error) {
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
	return &MarketplacePaginatedCategories{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *marketplaceCategoryService) GetCategoryByID(id uint) (*adminDomain.Category, error) {
	row, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMarketplaceCategoryNotFound
	}
	return row, nil
}

func (s *marketplaceCategoryService) GetSubcategories(page, limit int, parentID *uint) (*MarketplacePaginatedCategories, error) {
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
	return &MarketplacePaginatedCategories{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *marketplaceCategoryService) GetSubcategoryByID(id uint) (*adminDomain.Category, error) {
	row, err := s.repo.GetSubcategoryByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrMarketplaceSubcategoryNotFound
	}
	return row, nil
}
