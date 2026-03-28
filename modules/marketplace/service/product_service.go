package service

import (
	contrDomain "backend/modules/contragents/domain"
	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/repository"
)

type MarketplaceProductFilter = repository.MarketplaceProductFilter

type MarketplaceProductService interface {
	ListApproved(filter MarketplaceProductFilter, page, limit int) (*domain.PaginatedProducts, error)
	GetApprovedByID(id uint) (*domain.ProductOutput, error)
}

type marketplaceProductService struct {
	repo repository.MarketplaceProductRepository
}

func NewMarketplaceProductService(repo repository.MarketplaceProductRepository) MarketplaceProductService {
	return &marketplaceProductService{repo: repo}
}

func ProductOutputsFromRows(repo repository.MarketplaceProductRepository, rows []contrDomain.Product) ([]domain.ProductOutput, error) {
	productIDs := make([]uint, 0, len(rows))
	contragentIDs := make([]uint, 0, len(rows))
	seenContragent := map[uint]struct{}{}
	for _, row := range rows {
		productIDs = append(productIDs, row.ID)
		if _, ok := seenContragent[row.ContragentID]; !ok {
			seenContragent[row.ContragentID] = struct{}{}
			contragentIDs = append(contragentIDs, row.ContragentID)
		}
	}
	imagesMap, err := repo.GetImages(productIDs)
	if err != nil {
		return nil, err
	}
	deliveryAreasMap, err := repo.GetDeliveryAreasByContragentIDs(contragentIDs)
	if err != nil {
		return nil, err
	}
	items := make([]domain.ProductOutput, 0, len(rows))
	for i := range rows {
		items = append(items, mapMarketplaceProductOutput(&rows[i], imagesMap, deliveryAreasMap))
	}
	return items, nil
}

func (s *marketplaceProductService) ListApproved(filter MarketplaceProductFilter, page, limit int) (*domain.PaginatedProducts, error) {
	rows, total, err := s.repo.ListApprovedProducts(filter, page, limit)
	if err != nil {
		return nil, err
	}

	items, err := ProductOutputsFromRows(s.repo, rows)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &domain.PaginatedProducts{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *marketplaceProductService) GetApprovedByID(id uint) (*domain.ProductOutput, error) {
	row, err := s.repo.GetApprovedProductByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	imagesMap, err := s.repo.GetImages([]uint{row.ID})
	if err != nil {
		return nil, err
	}
	deliveryAreasMap, err := s.repo.GetDeliveryAreasByContragentIDs([]uint{row.ContragentID})
	if err != nil {
		return nil, err
	}

	out := mapMarketplaceProductOutput(row, imagesMap, deliveryAreasMap)
	return &out, nil
}

func mapMarketplaceProductOutput(row *contrDomain.Product, imagesMap map[uint][]string, deliveryAreasMap map[uint]domain.DeliveryAreas) domain.ProductOutput {
	margin := row.Price - row.OriginalPrice
	if margin < 0 {
		margin = 0
	}
	return domain.ProductOutput{
		ID:               row.ID,
		ProductCode:      row.ProductCode,
		ContragentID:     row.ContragentID,
		Name:             row.Name,
		Description:      row.Description,
		Price:            row.Price,
		OriginalPrice:    row.OriginalPrice,
		Images:           imagesMap[row.ID],
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
		DeliveryAreas:    deliveryAreasMap[row.ContragentID],
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}
