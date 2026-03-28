package domain

import "time"

// NeighborhoodShop — maxalla do'koni (viloyat, tuman, MFY bilan bog'langan).
type NeighborhoodShop struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:200;not null" json:"name"`
	INN                  string    `gorm:"size:12;index" json:"inn,omitempty"`
	RegionID             uint      `gorm:"not null;index" json:"region_id"`
	DistrictID           uint      `gorm:"not null;index" json:"district_id"`
	MFYID                uint      `gorm:"column:mfy_id;not null;index" json:"mfy_id"`
	Phone                string    `gorm:"size:13;uniqueIndex;not null" json:"phone"`
	Logo                 string    `gorm:"type:text" json:"logo,omitempty"`
	Status               string    `gorm:"size:20;not null" json:"status"`
	Password             string    `gorm:"size:255" json:"-"`
	PasswordSetupAllowed bool      `gorm:"not null;default:true" json:"password_setup_allowed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (NeighborhoodShop) TableName() string {
	return "neighborhood_shops"
}
