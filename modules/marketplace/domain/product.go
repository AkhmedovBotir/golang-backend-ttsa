package domain

import "time"

type DeliveryAreas struct {
	RegionIDs   []uint `json:"region_ids"`
	DistrictIDs []uint `json:"district_ids"`
}

type ProductOutput struct {
	ID               uint          `json:"id"`
	ProductCode      uint64        `json:"product_code"`
	ContragentID     uint          `json:"contragent_id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Price            float64       `json:"price"`
	OriginalPrice    float64       `json:"original_price"`
	Images           []string      `json:"images"`
	CategoryID       uint          `json:"category_id"`
	SubcategoryID    uint          `json:"subcategory_id"`
	Quantity         float64       `json:"quantity"`
	Unit             string        `json:"unit"`
	UnitSize         string        `json:"unit_size"`
	Status           string        `json:"status"`
	KpiBonusPercent  float64       `json:"kpi_bonus_percent"`
	KpiBonusAmount   float64       `json:"kpi_bonus_amount"`
	ModerationStatus string        `json:"moderation_status"`
	RejectionReason  string        `json:"rejection_reason"`
	DeliveryAreas    DeliveryAreas `json:"delivery_areas"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type PaginatedProducts struct {
	Items      []ProductOutput `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}
