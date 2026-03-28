package service

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/domain"
	"backend/modules/contragents/repository"
)

var (
	ErrProductNotFound            = errors.New("mahsulot topilmadi")
	ErrProductNameRequired        = errors.New("name majburiy")
	ErrProductDescriptionRequired = errors.New("description majburiy")
	ErrProductPriceInvalid        = errors.New("price noto'g'ri")
	ErrProductOriginalInvalid     = errors.New("original_price noto'g'ri")
	ErrProductCategoryInvalid     = errors.New("category_id noto'g'ri")
	ErrProductSubcategoryInvalid  = errors.New("subcategory_id noto'g'ri")
	ErrProductCategoryRelation    = errors.New("subcategory tanlangan categoryga tegishli emas")
	ErrProductQuantityInvalid     = errors.New("quantity noto'g'ri")
	ErrProductUnitInvalid         = errors.New("unit noto'g'ri")
	ErrProductUnitSizeRequired    = errors.New("unit_size majburiy")
	ErrProductImagesInvalid       = errors.New("images 1 tadan 5 tagacha bo'lishi kerak")
	ErrProductImageBase64Invalid  = errors.New("image base64 formati noto'g'ri")
	ErrProductKPIInvalid          = errors.New("kpi_bonus_percent 0 dan 100 gacha bo'lishi kerak")
	ErrProductStatusInvalid       = errors.New("status noto'g'ri")
)

type ProductInput struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Price           float64  `json:"price"`
	OriginalPrice   float64  `json:"original_price"`
	Images          []string `json:"images"`
	CategoryID      uint     `json:"category_id"`
	SubcategoryID   uint     `json:"subcategory_id"`
	Quantity        float64  `json:"quantity"`
	Unit            string   `json:"unit"`
	UnitSize        string   `json:"unit_size"`
	Status          string   `json:"status"`
	KpiBonusPercent float64  `json:"kpi_bonus_percent"`
	RejectionReason string   `json:"rejection_reason"`
}

type PaginatedProducts struct {
	Items      []ProductOutput `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

type ProductOutput struct {
	ID               uint      `json:"id"`
	ProductCode      uint64    `json:"product_code"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Price            float64   `json:"price"`
	OriginalPrice    float64   `json:"original_price"`
	Images           []string  `json:"images"`
	CategoryID       uint      `json:"category_id"`
	SubcategoryID    uint      `json:"subcategory_id"`
	Quantity         float64   `json:"quantity"`
	Unit             string    `json:"unit"`
	UnitSize         string    `json:"unit_size"`
	Status           string    `json:"status"`
	KpiBonusPercent  float64   `json:"kpi_bonus_percent"`
	KpiBonusAmount   float64   `json:"kpi_bonus_amount"`
	ModerationStatus string    `json:"moderation_status"`
	RejectionReason  string    `json:"rejection_reason"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ContragentProductService interface {
	Create(contragentID uint, input ProductInput) (*ProductOutput, error)
	GetPaginated(contragentID uint, page, limit int) (*PaginatedProducts, error)
	GetByID(contragentID, id uint) (*ProductOutput, error)
	Update(contragentID, id uint, input ProductInput) (*ProductOutput, error)
	UpdateStatus(contragentID, id uint, status string) (*ProductOutput, error)
	Delete(contragentID, id uint) error
}

type contragentProductService struct {
	repo repository.ContragentProductRepository
}

func NewContragentProductService(repo repository.ContragentProductRepository) ContragentProductService {
	return &contragentProductService{repo: repo}
}

func (s *contragentProductService) Create(contragentID uint, input ProductInput) (*ProductOutput, error) {
	input.normalize()
	if err := s.validateInput(input, true); err != nil {
		return nil, err
	}
	if err := s.validateCategoryRelation(input.CategoryID, input.SubcategoryID); err != nil {
		return nil, err
	}
	code, err := s.repo.NextProductCode()
	if err != nil {
		return nil, err
	}
	row := &domain.Product{
		ContragentID:     contragentID,
		ProductCode:      code,
		Name:             input.Name,
		Description:      input.Description,
		Price:            input.Price,
		OriginalPrice:    input.OriginalPrice,
		CategoryID:       input.CategoryID,
		SubcategoryID:    input.SubcategoryID,
		Quantity:         input.Quantity,
		Unit:             input.Unit,
		UnitSize:         input.UnitSize,
		Status:           input.Status,
		KpiBonusPercent:  input.KpiBonusPercent,
		ModerationStatus: domain.ProductModerationPending,
		RejectionReason:  "",
	}
	if err = s.repo.Create(row, input.Images); err != nil {
		return nil, err
	}
	return s.GetByID(contragentID, row.ID)
}

func (s *contragentProductService) GetPaginated(contragentID uint, page, limit int) (*PaginatedProducts, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	rows, total, err := s.repo.GetPaginatedByContragent(contragentID, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]ProductOutput, 0, len(rows))
	for _, row := range rows {
		images, err := s.repo.GetImages(row.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, mapProductOutput(&row, images))
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &PaginatedProducts{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *contragentProductService) GetByID(contragentID, id uint) (*ProductOutput, error) {
	row, err := s.repo.GetByIDAndContragent(id, contragentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrProductNotFound
	}
	images, err := s.repo.GetImages(row.ID)
	if err != nil {
		return nil, err
	}
	out := mapProductOutput(row, images)
	return &out, nil
}

func (s *contragentProductService) Update(contragentID, id uint, input ProductInput) (*ProductOutput, error) {
	input.normalize()
	if err := s.validateInput(input, true); err != nil {
		return nil, err
	}
	if err := s.validateCategoryRelation(input.CategoryID, input.SubcategoryID); err != nil {
		return nil, err
	}
	row, err := s.repo.GetByIDAndContragent(id, contragentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrProductNotFound
	}
	row.Name = input.Name
	row.Description = input.Description
	row.Price = input.Price
	row.OriginalPrice = input.OriginalPrice
	row.CategoryID = input.CategoryID
	row.SubcategoryID = input.SubcategoryID
	row.Quantity = input.Quantity
	row.Unit = input.Unit
	row.UnitSize = input.UnitSize
	row.Status = input.Status
	row.KpiBonusPercent = input.KpiBonusPercent
	// Kontragent update qilganda qayta moderatsiya jarayoniga tushadi.
	row.ModerationStatus = domain.ProductModerationPending
	row.RejectionReason = ""
	if err = s.repo.Update(row, input.Images); err != nil {
		return nil, err
	}
	return s.GetByID(contragentID, row.ID)
}

func (s *contragentProductService) UpdateStatus(contragentID, id uint, status string) (*ProductOutput, error) {
	row, err := s.repo.GetByIDAndContragent(id, contragentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrProductNotFound
	}
	st, err := normalizeProductStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.Update(row, nil); err != nil {
		return nil, err
	}
	return s.GetByID(contragentID, row.ID)
}

func (s *contragentProductService) Delete(contragentID, id uint) error {
	row, err := s.repo.GetByIDAndContragent(id, contragentID)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrProductNotFound
	}
	return s.repo.Delete(row.ID)
}

func (s *contragentProductService) validateInput(input ProductInput, withImages bool) error {
	if input.Name == "" {
		return ErrProductNameRequired
	}
	if input.Description == "" {
		return ErrProductDescriptionRequired
	}
	if input.Price <= 0 {
		return ErrProductPriceInvalid
	}
	if input.OriginalPrice <= 0 {
		return ErrProductOriginalInvalid
	}
	if input.CategoryID == 0 {
		return ErrProductCategoryInvalid
	}
	if input.SubcategoryID == 0 {
		return ErrProductSubcategoryInvalid
	}
	if input.Quantity <= 0 {
		return ErrProductQuantityInvalid
	}
	if input.Unit != domain.ProductUnitDona && input.Unit != domain.ProductUnitLitr && input.Unit != domain.ProductUnitKg {
		return ErrProductUnitInvalid
	}
	if input.UnitSize == "" {
		return ErrProductUnitSizeRequired
	}
	if input.KpiBonusPercent < 0 || input.KpiBonusPercent > 100 {
		return ErrProductKPIInvalid
	}
	if _, err := normalizeProductStatus(input.Status); err != nil {
		return err
	}
	if withImages {
		if len(input.Images) < 1 || len(input.Images) > 5 {
			return ErrProductImagesInvalid
		}
		for _, img := range input.Images {
			if !isValidBase64ProductImage(img) {
				return ErrProductImageBase64Invalid
			}
		}
	}
	return nil
}

func (s *contragentProductService) validateCategoryRelation(categoryID, subcategoryID uint) error {
	cat, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	if cat == nil || cat.ParentID != nil {
		return ErrProductCategoryInvalid
	}
	sub, err := s.repo.GetSubcategoryByID(subcategoryID)
	if err != nil {
		return err
	}
	if sub == nil || sub.ParentID == nil {
		return ErrProductSubcategoryInvalid
	}
	if *sub.ParentID != categoryID {
		return ErrProductCategoryRelation
	}
	return nil
}

func mapProductOutput(row *domain.Product, images []string) ProductOutput {
	margin := row.Price - row.OriginalPrice
	if margin < 0 {
		margin = 0
	}
	return ProductOutput{
		ID:               row.ID,
		ProductCode:      row.ProductCode,
		Name:             row.Name,
		Description:      row.Description,
		Price:            row.Price,
		OriginalPrice:    row.OriginalPrice,
		Images:           images,
		CategoryID:       row.CategoryID,
		SubcategoryID:    row.SubcategoryID,
		Quantity:         row.Quantity,
		Unit:             row.Unit,
		UnitSize:         row.UnitSize,
		Status:           row.Status,
		KpiBonusPercent:  row.KpiBonusPercent,
		KpiBonusAmount:   margin * row.KpiBonusPercent / 100,
		ModerationStatus: row.ModerationStatus,
		RejectionReason:  row.RejectionReason,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func (i *ProductInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Description = strings.TrimSpace(i.Description)
	i.Unit = strings.TrimSpace(i.Unit)
	i.UnitSize = strings.TrimSpace(i.UnitSize)
	i.Status = strings.TrimSpace(i.Status)
	i.RejectionReason = strings.TrimSpace(i.RejectionReason)
	for idx := range i.Images {
		i.Images[idx] = strings.TrimSpace(i.Images[idx])
	}
}

func isValidBase64ProductImage(raw string) bool {
	payload := raw
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return false
		}
		payload = parts[1]
	}
	if payload == "" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(payload)
	return err == nil
}

func normalizeProductStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	if status != adminDomain.StatusActive && status != adminDomain.StatusInactive {
		return "", ErrProductStatusInvalid
	}
	return status, nil
}
