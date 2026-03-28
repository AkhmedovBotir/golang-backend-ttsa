package service

import (
	"sort"
	"strings"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/repository"
)

type ContragentBrowseInclude struct {
	Products   bool
	Categories bool
}

func ParseContragentInclude(raw string) ContragentBrowseInclude {
	var inc ContragentBrowseInclude
	for _, part := range strings.Split(raw, ",") {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "products", "product":
			inc.Products = true
		case "categories", "category", "subcategories", "subcategory":
			inc.Categories = true
		}
	}
	return inc
}

type MarketplaceContragentBrowseService interface {
	List(query *string, page, limit, nestedLimit int, inc ContragentBrowseInclude) (*domain.PaginatedContragents, error)
}

type marketplaceContragentBrowseService struct {
	browseRepo  repository.MarketplaceContragentBrowseRepository
	productRepo repository.MarketplaceProductRepository
}

func NewMarketplaceContragentBrowseService(
	browseRepo repository.MarketplaceContragentBrowseRepository,
	productRepo repository.MarketplaceProductRepository,
) MarketplaceContragentBrowseService {
	return &marketplaceContragentBrowseService{browseRepo: browseRepo, productRepo: productRepo}
}

func (s *marketplaceContragentBrowseService) List(query *string, page, limit, nestedLimit int, inc ContragentBrowseInclude) (*domain.PaginatedContragents, error) {
	if nestedLimit < 1 {
		nestedLimit = 30
	}
	if nestedLimit > 100 {
		nestedLimit = 100
	}

	rows, total, err := s.browseRepo.ListActiveContragents(query, page, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(rows))
	for i := range rows {
		ids = append(ids, rows[i].ID)
	}
	deliveryMap, err := s.productRepo.GetDeliveryAreasByContragentIDs(ids)
	if err != nil {
		return nil, err
	}

	var pairs []repository.CatSubPairForContragent
	if inc.Categories && len(ids) > 0 {
		pairs, err = s.browseRepo.GetDistinctCategoryPairsForContragents(ids)
		if err != nil {
			return nil, err
		}
	}

	catIDSet := map[uint]struct{}{}
	for _, p := range pairs {
		catIDSet[p.CategoryID] = struct{}{}
		catIDSet[p.SubcategoryID] = struct{}{}
	}
	catIDs := make([]uint, 0, len(catIDSet))
	for id := range catIDSet {
		catIDs = append(catIDs, id)
	}
	catByID, err := s.browseRepo.GetCategoriesByIDs(catIDs)
	if err != nil {
		return nil, err
	}

	pairsByAgent := map[uint][]repository.CatSubPairForContragent{}
	for _, p := range pairs {
		pairsByAgent[p.ContragentID] = append(pairsByAgent[p.ContragentID], p)
	}

	items := make([]domain.ContragentBrowseItem, 0, len(rows))
	for i := range rows {
		a := &rows[i]
		item := domain.ContragentBrowseItem{
			ID: a.ID, Name: a.Name, INN: a.INN, Phone: a.Phone, Logo: a.Logo,
			RegionID: a.RegionID, DistrictID: a.DistrictID, MFYID: a.MFYID,
			ActivityTypeID: a.ActivityTypeID, Status: a.Status,
			DeliveryAreas: deliveryMap[a.ID],
		}
		if inc.Products {
			prows, err := s.browseRepo.ListApprovedProductsForContragent(a.ID, nestedLimit)
			if err != nil {
				return nil, err
			}
			item.Products, err = ProductOutputsFromRows(s.productRepo, prows)
			if err != nil {
				return nil, err
			}
		}
		if inc.Categories {
			item.CategoryBranches = buildCategoryBranches(pairsByAgent[a.ID], catByID)
		}
		items = append(items, item)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &domain.PaginatedContragents{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func buildCategoryBranches(pairs []repository.CatSubPairForContragent, catByID map[uint]adminDomain.Category) []domain.ContragentCatBranch {
	catToSubs := map[uint]map[uint]struct{}{}
	for _, p := range pairs {
		if catToSubs[p.CategoryID] == nil {
			catToSubs[p.CategoryID] = map[uint]struct{}{}
		}
		catToSubs[p.CategoryID][p.SubcategoryID] = struct{}{}
	}

	branches := make([]domain.ContragentCatBranch, 0, len(catToSubs))
	for catID, subSet := range catToSubs {
		main, ok := catByID[catID]
		if !ok || main.ParentID != nil {
			continue
		}
		subs := make([]domain.CategorySearchHit, 0, len(subSet))
		for subID := range subSet {
			sub, ok := catByID[subID]
			if !ok || sub.ParentID == nil {
				continue
			}
			if *sub.ParentID != catID {
				continue
			}
			subs = append(subs, categorySearchHitFromAdmin(&sub))
		}
		sort.Slice(subs, func(i, j int) bool { return subs[i].ID < subs[j].ID })
		branches = append(branches, domain.ContragentCatBranch{
			Category:      categorySearchHitFromAdmin(&main),
			Subcategories: subs,
		})
	}
	sort.Slice(branches, func(i, j int) bool { return branches[i].Category.ID < branches[j].Category.ID })
	return branches
}
