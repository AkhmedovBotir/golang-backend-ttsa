package service

import (
	"errors"
	"strings"

	"backend/modules/admin/domain"
	"backend/modules/admin/repository"
)

var (
	ErrCategoryNotFound      = errors.New("kategoriya topilmadi")
	ErrSubcategoryNotFound   = errors.New("subkategoriya topilmadi")
	ErrCategorySlugExists    = errors.New("slug allaqachon mavjud")
	ErrCategoryNameRequired  = errors.New("nom majburiy")
	ErrCategorySlugRequired  = errors.New("slug majburiy")
	ErrCategoryParentInvalid = errors.New("parent_id noto'g'ri")
)

type CategoryInput struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Image    string `json:"image"`
	Censored *bool  `json:"censored"`
	Status   string `json:"status"`
}

type SubcategoryInput struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Image    string `json:"image"`
	Censored *bool  `json:"censored"`
	ParentID uint   `json:"parent_id"`
	Status   string `json:"status"`
}

type PaginatedCategories struct {
	Items      []domain.Category `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

type CategoryService interface {
	CreateCategory(input CategoryInput) (*domain.Category, error)
	GetCategories(page, limit int) (*PaginatedCategories, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(id uint, input CategoryInput) (*domain.Category, error)
	UpdateCategoryStatus(id uint, status string) (*domain.Category, error)
	DeleteCategory(id uint) error

	CreateSubcategory(input SubcategoryInput) (*domain.Category, error)
	GetSubcategories(page, limit int, parentID *uint) (*PaginatedCategories, error)
	GetSubcategoryByID(id uint) (*domain.Category, error)
	UpdateSubcategory(id uint, input SubcategoryInput) (*domain.Category, error)
	UpdateSubcategoryStatus(id uint, status string) (*domain.Category, error)
	DeleteSubcategory(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(input CategoryInput) (*domain.Category, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrCategoryNameRequired
	}
	if input.Slug == "" {
		return nil, ErrCategorySlugRequired
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsBySlug(input.Slug, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrCategorySlugExists
	}
	censored := false
	if input.Censored != nil {
		censored = *input.Censored
	}
	row := &domain.Category{
		Name:     input.Name,
		Slug:     input.Slug,
		Image:    input.Image,
		Censored: censored,
		ParentID: nil,
		Status:   status,
	}
	return row, s.repo.Create(row)
}

func (s *categoryService) GetCategories(page, limit int) (*PaginatedCategories, error) {
	return s.getPaginated(page, limit, nil, true)
}

func (s *categoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil || row.ParentID != nil {
		return nil, ErrCategoryNotFound
	}
	return row, nil
}

func (s *categoryService) UpdateCategory(id uint, input CategoryInput) (*domain.Category, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrCategoryNameRequired
	}
	if input.Slug == "" {
		return nil, ErrCategorySlugRequired
	}
	row, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsBySlug(input.Slug, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrCategorySlugExists
	}
	censored := row.Censored
	if input.Censored != nil {
		censored = *input.Censored
	}
	row.Name = input.Name
	row.Slug = input.Slug
	row.Image = input.Image
	row.Censored = censored
	row.Status = status
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	// Asosiy kategoriya statusi o'zgarsa subcategorylar ham shu statusga o'tadi.
	if status == domain.StatusInactive || status == domain.StatusActive {
		if err = s.repo.UpdateChildrenStatusByParentID(row.ID, status); err != nil {
			return nil, err
		}
	}
	return row, nil
}

func (s *categoryService) UpdateCategoryStatus(id uint, status string) (*domain.Category, error) {
	row, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	st, err := normalizeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.Update(row); err != nil {
		return nil, err
	}
	// Asosiy kategoriya statusi o'zgarsa subcategorylar ham shu statusga o'tadi.
	if st == domain.StatusInactive || st == domain.StatusActive {
		if err = s.repo.UpdateChildrenStatusByParentID(row.ID, st); err != nil {
			return nil, err
		}
	}
	return row, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	row, err := s.GetCategoryByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(row.ID)
}

func (s *categoryService) CreateSubcategory(input SubcategoryInput) (*domain.Category, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrCategoryNameRequired
	}
	if input.Slug == "" {
		return nil, ErrCategorySlugRequired
	}
	if input.ParentID == 0 {
		return nil, ErrCategoryParentInvalid
	}
	if ok, err := s.repo.ParentExists(input.ParentID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrCategoryParentInvalid
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsBySlug(input.Slug, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrCategorySlugExists
	}
	censored := false
	if input.Censored != nil {
		censored = *input.Censored
	}
	parentID := input.ParentID
	row := &domain.Category{
		Name:     input.Name,
		Slug:     input.Slug,
		Image:    input.Image,
		Censored: censored,
		ParentID: &parentID,
		Status:   status,
	}
	return row, s.repo.Create(row)
}

func (s *categoryService) GetSubcategories(page, limit int, parentID *uint) (*PaginatedCategories, error) {
	return s.getPaginated(page, limit, parentID, false)
}

func (s *categoryService) GetSubcategoryByID(id uint) (*domain.Category, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil || row.ParentID == nil {
		return nil, ErrSubcategoryNotFound
	}
	return row, nil
}

func (s *categoryService) UpdateSubcategory(id uint, input SubcategoryInput) (*domain.Category, error) {
	input.normalize()
	if input.Name == "" {
		return nil, ErrCategoryNameRequired
	}
	if input.Slug == "" {
		return nil, ErrCategorySlugRequired
	}
	if input.ParentID == 0 {
		return nil, ErrCategoryParentInvalid
	}
	if ok, err := s.repo.ParentExists(input.ParentID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrCategoryParentInvalid
	}
	row, err := s.GetSubcategoryByID(id)
	if err != nil {
		return nil, err
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if exists, err := s.repo.ExistsBySlug(input.Slug, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrCategorySlugExists
	}
	censored := row.Censored
	if input.Censored != nil {
		censored = *input.Censored
	}
	parentID := input.ParentID
	row.Name = input.Name
	row.Slug = input.Slug
	row.Image = input.Image
	row.Censored = censored
	row.ParentID = &parentID
	row.Status = status
	return row, s.repo.Update(row)
}

func (s *categoryService) UpdateSubcategoryStatus(id uint, status string) (*domain.Category, error) {
	row, err := s.GetSubcategoryByID(id)
	if err != nil {
		return nil, err
	}
	st, err := normalizeStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	return row, s.repo.Update(row)
}

func (s *categoryService) DeleteSubcategory(id uint) error {
	row, err := s.GetSubcategoryByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(row.ID)
}

func (s *categoryService) getPaginated(page, limit int, parentID *uint, parentOnly bool) (*PaginatedCategories, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	var (
		items []domain.Category
		total int64
		err   error
	)
	if parentOnly {
		items, total, err = s.repo.GetPaginatedCategories(page, limit)
	} else {
		items, total, err = s.repo.GetPaginatedSubcategories(page, limit, parentID)
	}
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &PaginatedCategories{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (i *CategoryInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Slug = strings.TrimSpace(i.Slug)
	i.Image = strings.TrimSpace(i.Image)
	i.Status = strings.TrimSpace(i.Status)
}

func (i *SubcategoryInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Slug = strings.TrimSpace(i.Slug)
	i.Image = strings.TrimSpace(i.Image)
	i.Status = strings.TrimSpace(i.Status)
}
