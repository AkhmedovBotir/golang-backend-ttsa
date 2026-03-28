package service

import (
	"strings"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/repository"
)

type SearchIncludeMask uint8

const (
	IncludeSearchProducts SearchIncludeMask = 1 << iota
	IncludeSearchCategories
	IncludeSearchSubcategories
	IncludeSearchContragents
)

const includeSearchAll = IncludeSearchProducts | IncludeSearchCategories | IncludeSearchSubcategories | IncludeSearchContragents

func ParseSearchTypes(raw string) SearchIncludeMask {
	if strings.TrimSpace(raw) == "" {
		return includeSearchAll
	}
	var m SearchIncludeMask
	for _, part := range strings.Split(raw, ",") {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "products", "product":
			m |= IncludeSearchProducts
		case "categories", "category":
			m |= IncludeSearchCategories
		case "subcategories", "subcategory":
			m |= IncludeSearchSubcategories
		case "contragents", "contragent":
			m |= IncludeSearchContragents
		}
	}
	if m == 0 {
		return includeSearchAll
	}
	return m
}

type MarketplaceSearchService interface {
	UnifiedSearch(query string, limitPerType int, mask SearchIncludeMask) (*domain.UnifiedSearchResponse, error)
}

type marketplaceSearchService struct {
	searchRepo  repository.MarketplaceSearchRepository
	productRepo repository.MarketplaceProductRepository
}

func NewMarketplaceSearchService(
	searchRepo repository.MarketplaceSearchRepository,
	productRepo repository.MarketplaceProductRepository,
) MarketplaceSearchService {
	return &marketplaceSearchService{searchRepo: searchRepo, productRepo: productRepo}
}

func (s *marketplaceSearchService) UnifiedSearch(query string, limitPerType int, mask SearchIncludeMask) (*domain.UnifiedSearchResponse, error) {
	q := strings.TrimSpace(query)
	if limitPerType < 1 {
		limitPerType = 10
	}
	if limitPerType > 50 {
		limitPerType = 50
	}

	out := &domain.UnifiedSearchResponse{
		Query:         q,
		LimitPerType:  limitPerType,
		Products:      []domain.ProductOutput{},
		Categories:    []domain.CategorySearchHit{},
		Subcategories: []domain.CategorySearchHit{},
		Contragents:   []domain.ContragentSearchHit{},
	}
	if q == "" {
		return out, nil
	}

	if mask&IncludeSearchProducts != 0 {
		prows, err := s.searchRepo.SearchProducts(q, limitPerType)
		if err != nil {
			return nil, err
		}
		items, err := ProductOutputsFromRows(s.productRepo, prows)
		if err != nil {
			return nil, err
		}
		out.Products = items
	}

	if mask&IncludeSearchCategories != 0 {
		crows, err := s.searchRepo.SearchMainCategories(q, limitPerType)
		if err != nil {
			return nil, err
		}
		for i := range crows {
			out.Categories = append(out.Categories, categorySearchHitFromAdmin(&crows[i]))
		}
	}

	if mask&IncludeSearchSubcategories != 0 {
		srows, err := s.searchRepo.SearchSubcategories(q, limitPerType)
		if err != nil {
			return nil, err
		}
		for i := range srows {
			out.Subcategories = append(out.Subcategories, categorySearchHitFromAdmin(&srows[i]))
		}
	}

	if mask&IncludeSearchContragents != 0 {
		agrows, err := s.searchRepo.SearchContragents(q, limitPerType)
		if err != nil {
			return nil, err
		}
		agentIDs := make([]uint, 0, len(agrows))
		for i := range agrows {
			agentIDs = append(agentIDs, agrows[i].ID)
		}
		deliveryByAgent, err := s.productRepo.GetDeliveryAreasByContragentIDs(agentIDs)
		if err != nil {
			return nil, err
		}
		for i := range agrows {
			a := &agrows[i]
			out.Contragents = append(out.Contragents, domain.ContragentSearchHit{
				ID: a.ID, Name: a.Name, INN: a.INN, Phone: a.Phone, Logo: a.Logo,
				RegionID: a.RegionID, DistrictID: a.DistrictID, MFYID: a.MFYID,
				DeliveryAreas: deliveryByAgent[a.ID],
			})
		}
	}

	return out, nil
}

func categorySearchHitFromAdmin(c *adminDomain.Category) domain.CategorySearchHit {
	return domain.CategorySearchHit{
		ID: c.ID, Name: c.Name, Slug: c.Slug, Image: c.Image, ParentID: c.ParentID,
	}
}
