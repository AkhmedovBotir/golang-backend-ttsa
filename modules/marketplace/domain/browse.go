package domain

// UnifiedSearchResponse — bitta so'rov bo'yicha mahsulot, kategoriya, subkategoriya va kontragentlar.
type UnifiedSearchResponse struct {
	Query         string                `json:"query"`
	LimitPerType  int                   `json:"limit_per_type"`
	Products      []ProductOutput       `json:"products"`
	Categories    []CategorySearchHit   `json:"categories"`
	Subcategories []CategorySearchHit   `json:"subcategories"`
	Contragents   []ContragentSearchHit `json:"contragents"`
}

type CategorySearchHit struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Image    string `json:"image"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type ContragentSearchHit struct {
	ID            uint          `json:"id"`
	Name          string        `json:"name"`
	INN           string        `json:"inn"`
	Phone         string        `json:"phone"`
	Logo          string        `json:"logo,omitempty"`
	RegionID      uint          `json:"region_id"`
	DistrictID    uint          `json:"district_id"`
	MFYID         uint          `json:"mfy_id"`
	DeliveryAreas DeliveryAreas `json:"delivery_areas"`
}

// ContragentBrowseItem — marketplace kontragenti joylashuv va yetkazib berish hududi bilan.
type ContragentBrowseItem struct {
	ID               uint                  `json:"id"`
	Name             string                `json:"name"`
	INN              string                `json:"inn"`
	Phone            string                `json:"phone"`
	Logo             string                `json:"logo,omitempty"`
	RegionID         uint                  `json:"region_id"`
	DistrictID       uint                  `json:"district_id"`
	MFYID            uint                  `json:"mfy_id"`
	ActivityTypeID   uint                  `json:"activity_type_id"`
	Status           string                `json:"status"`
	DeliveryAreas    DeliveryAreas         `json:"delivery_areas"`
	Products         []ProductOutput       `json:"products,omitempty"`
	CategoryBranches []ContragentCatBranch `json:"category_branches,omitempty"`
}

type ContragentCatBranch struct {
	Category      CategorySearchHit   `json:"category"`
	Subcategories []CategorySearchHit `json:"subcategories"`
}

type PaginatedContragents struct {
	Items      []ContragentBrowseItem `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}
