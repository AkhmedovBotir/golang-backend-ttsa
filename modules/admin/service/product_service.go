package service

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/admin/repository"
	contrDomain "backend/modules/contragents/domain"
)

var (
	ErrAdminProductNotFound            = errors.New("mahsulot topilmadi")
	ErrAdminProductNameRequired        = errors.New("name majburiy")
	ErrAdminProductDescriptionRequired = errors.New("description majburiy")
	ErrAdminProductPriceInvalid        = errors.New("price noto'g'ri")
	ErrAdminProductOriginalInvalid     = errors.New("original_price noto'g'ri")
	ErrAdminProductCategoryInvalid     = errors.New("category_id noto'g'ri")
	ErrAdminProductSubcategoryInvalid  = errors.New("subcategory_id noto'g'ri")
	ErrAdminProductCategoryRelation    = errors.New("subcategory tanlangan categoryga tegishli emas")
	ErrAdminProductQuantityInvalid     = errors.New("quantity noto'g'ri")
	ErrAdminProductUnitInvalid         = errors.New("unit noto'g'ri")
	ErrAdminProductUnitSizeRequired    = errors.New("unit_size majburiy")
	ErrAdminProductImagesInvalid       = errors.New("images 1 tadan 5 tagacha bo'lishi kerak")
	ErrAdminProductImageBase64Invalid  = errors.New("image base64 formati noto'g'ri")
	ErrAdminProductKPIInvalid          = errors.New("kpi_bonus_percent 0 dan 100 gacha bo'lishi kerak")
	ErrAdminProductStatusInvalid       = errors.New("status noto'g'ri")
	ErrAdminProductContragentInvalid   = errors.New("contragent_id noto'g'ri")
	ErrAdminProductModerationInvalid   = errors.New("moderation_status noto'g'ri")
	ErrAdminRejectReasonRequired       = errors.New("rejection_reason majburiy")
)

type AdminProductInput struct {
	ContragentID    uint     `json:"contragent_id"`
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
}

type AdminProductOutput struct {
	ID               uint      `json:"id"`
	ContragentID     uint      `json:"contragent_id"`
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

type AdminPaginatedProducts struct {
	Items      []AdminProductOutput `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

type AdminProductService interface {
	Create(input AdminProductInput) (*AdminProductOutput, error)
	GetPaginated(page, limit int, contragentID *uint, moderationStatus *string) (*AdminPaginatedProducts, error)
	GetByID(id uint) (*AdminProductOutput, error)
	Update(id uint, input AdminProductInput) (*AdminProductOutput, error)
	UpdateStatus(id uint, status string) (*AdminProductOutput, error)
	Delete(id uint) error
	Approve(id uint) (*AdminProductOutput, error)
	Reject(id uint, reason string) (*AdminProductOutput, error)
}

type adminProductService struct {
	repo repository.AdminProductRepository
}

func NewAdminProductService(repo repository.AdminProductRepository) AdminProductService {
	return &adminProductService{repo: repo}
}

func (s *adminProductService) Create(input AdminProductInput) (*AdminProductOutput, error) {
	input.normalize()
	if err := s.validateInput(input, true); err != nil {
		return nil, err
	}
	if err := s.validateRelations(input.ContragentID, input.CategoryID, input.SubcategoryID); err != nil {
		return nil, err
	}
	code, err := s.repo.NextProductCode()
	if err != nil {
		return nil, err
	}
	row := &contrDomain.Product{
		ContragentID:     input.ContragentID,
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
		ModerationStatus: contrDomain.ProductModerationPending,
		RejectionReason:  "",
	}
	if err = s.repo.Create(row, input.Images); err != nil {
		return nil, err
	}
	return s.GetByID(row.ID)
}

func (s *adminProductService) GetPaginated(page, limit int, contragentID *uint, moderationStatus *string) (*AdminPaginatedProducts, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	rows, total, err := s.repo.GetPaginated(page, limit, contragentID, moderationStatus)
	if err != nil {
		return nil, err
	}
	items := make([]AdminProductOutput, 0, len(rows))
	for _, row := range rows {
		images, err := s.repo.GetImages(row.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, mapAdminProductOutput(&row, images))
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return &AdminPaginatedProducts{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *adminProductService) GetByID(id uint) (*AdminProductOutput, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAdminProductNotFound
	}
	images, err := s.repo.GetImages(row.ID)
	if err != nil {
		return nil, err
	}
	out := mapAdminProductOutput(row, images)
	return &out, nil
}

func (s *adminProductService) Update(id uint, input AdminProductInput) (*AdminProductOutput, error) {
	input.normalize()
	if err := s.validateInput(input, true); err != nil {
		return nil, err
	}
	if err := s.validateRelations(input.ContragentID, input.CategoryID, input.SubcategoryID); err != nil {
		return nil, err
	}
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAdminProductNotFound
	}
	row.ContragentID = input.ContragentID
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
	if err = s.repo.Update(row, input.Images); err != nil {
		return nil, err
	}
	return s.GetByID(row.ID)
}

func (s *adminProductService) UpdateStatus(id uint, status string) (*AdminProductOutput, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAdminProductNotFound
	}
	st, err := normalizeAdminProductStatus(status)
	if err != nil {
		return nil, err
	}
	row.Status = st
	if err = s.repo.Update(row, nil); err != nil {
		return nil, err
	}
	return s.GetByID(row.ID)
}

func (s *adminProductService) Delete(id uint) error {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrAdminProductNotFound
	}
	return s.repo.Delete(row.ID)
}

func (s *adminProductService) Approve(id uint) (*AdminProductOutput, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAdminProductNotFound
	}
	row.ModerationStatus = contrDomain.ProductModerationApproved
	row.RejectionReason = ""
	if err = s.repo.Update(row, nil); err != nil {
		return nil, err
	}
	return s.GetByID(row.ID)
}

func (s *adminProductService) Reject(id uint, reason string) (*AdminProductOutput, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrAdminProductNotFound
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrAdminRejectReasonRequired
	}
	row.ModerationStatus = contrDomain.ProductModerationRejected
	row.RejectionReason = reason
	if err = s.repo.Update(row, nil); err != nil {
		return nil, err
	}
	return s.GetByID(row.ID)
}

func (s *adminProductService) validateInput(input AdminProductInput, withImages bool) error {
	if input.ContragentID == 0 {
		return ErrAdminProductContragentInvalid
	}
	if input.Name == "" {
		return ErrAdminProductNameRequired
	}
	if input.Description == "" {
		return ErrAdminProductDescriptionRequired
	}
	if input.Price <= 0 {
		return ErrAdminProductPriceInvalid
	}
	if input.OriginalPrice <= 0 {
		return ErrAdminProductOriginalInvalid
	}
	if input.CategoryID == 0 {
		return ErrAdminProductCategoryInvalid
	}
	if input.SubcategoryID == 0 {
		return ErrAdminProductSubcategoryInvalid
	}
	if input.Quantity <= 0 {
		return ErrAdminProductQuantityInvalid
	}
	if input.Unit != contrDomain.ProductUnitDona && input.Unit != contrDomain.ProductUnitLitr && input.Unit != contrDomain.ProductUnitKg {
		return ErrAdminProductUnitInvalid
	}
	if input.UnitSize == "" {
		return ErrAdminProductUnitSizeRequired
	}
	if input.KpiBonusPercent < 0 || input.KpiBonusPercent > 100 {
		return ErrAdminProductKPIInvalid
	}
	if _, err := normalizeAdminProductStatus(input.Status); err != nil {
		return err
	}
	if withImages {
		if len(input.Images) < 1 || len(input.Images) > 5 {
			return ErrAdminProductImagesInvalid
		}
		for _, img := range input.Images {
			if !isValidAdminProductImage(img) {
				return ErrAdminProductImageBase64Invalid
			}
		}
	}
	return nil
}

func (s *adminProductService) validateRelations(contragentID, categoryID, subcategoryID uint) error {
	ok, err := s.repo.ContragentExists(contragentID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAdminProductContragentInvalid
	}
	cat, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	if cat == nil || cat.ParentID != nil {
		return ErrAdminProductCategoryInvalid
	}
	sub, err := s.repo.GetSubcategoryByID(subcategoryID)
	if err != nil {
		return err
	}
	if sub == nil || sub.ParentID == nil {
		return ErrAdminProductSubcategoryInvalid
	}
	if *sub.ParentID != categoryID {
		return ErrAdminProductCategoryRelation
	}
	return nil
}

func mapAdminProductOutput(row *contrDomain.Product, images []string) AdminProductOutput {
	margin := row.Price - row.OriginalPrice
	if margin < 0 {
		margin = 0
	}
	return AdminProductOutput{
		ID:               row.ID,
		ContragentID:     row.ContragentID,
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

func normalizeAdminProductStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	if status != adminDomain.StatusActive && status != adminDomain.StatusInactive {
		return "", ErrAdminProductStatusInvalid
	}
	return status, nil
}

func isValidAdminProductImage(raw string) bool {
	payload := strings.TrimSpace(raw)
	if strings.HasPrefix(payload, "data:") {
		parts := strings.SplitN(payload, ",", 2)
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

func (i *AdminProductInput) normalize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Description = strings.TrimSpace(i.Description)
	i.Unit = strings.TrimSpace(i.Unit)
	i.UnitSize = strings.TrimSpace(i.UnitSize)
	i.Status = strings.TrimSpace(i.Status)
	for idx := range i.Images {
		i.Images[idx] = strings.TrimSpace(i.Images[idx])
	}
}
