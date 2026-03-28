package domain

import "time"

const (
	ProductUnitDona = "dona"
	ProductUnitLitr = "litr"
	ProductUnitKg   = "kg"

	ProductModerationPending  = "pending"
	ProductModerationApproved = "approved"
	ProductModerationRejected = "rejected"
)

type Product struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ContragentID     uint      `gorm:"index;not null" json:"contragent_id"`
	ProductCode      uint64    `gorm:"index;not null" json:"product_code"`
	Name             string    `gorm:"size:255;not null" json:"name"`
	Description      string    `gorm:"type:text;not null" json:"description"`
	Price            float64   `gorm:"not null" json:"price"`
	OriginalPrice    float64   `gorm:"column:original_price;not null" json:"original_price"`
	CategoryID       uint      `gorm:"index;not null" json:"category_id"`
	SubcategoryID    uint      `gorm:"index;not null" json:"subcategory_id"`
	Quantity         float64   `gorm:"not null" json:"quantity"`
	Unit             string    `gorm:"size:16;not null" json:"unit"`
	UnitSize         string    `gorm:"size:120;not null" json:"unit_size"`
	Status           string    `gorm:"size:20;not null" json:"status"`
	KpiBonusPercent  float64   `gorm:"column:kpi_bonus_percent;not null;default:0" json:"kpi_bonus_percent"`
	ModerationStatus string    `gorm:"size:20;not null;default:pending" json:"moderation_status"`
	RejectionReason  string    `gorm:"type:text" json:"rejection_reason"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Image     string    `gorm:"type:text;not null" json:"image"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
