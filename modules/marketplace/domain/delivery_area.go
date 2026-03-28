package domain

import "time"

type DeliveryArea struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Name       string    `gorm:"size:120;not null" json:"name"`
	RegionID   uint      `gorm:"not null;index" json:"region_id"`
	DistrictID uint      `gorm:"not null;index" json:"district_id"`
	MFYID      uint      `gorm:"column:mfy_id;not null;index" json:"mfy_id"`
	IsDefault  bool      `gorm:"not null;default:false;index" json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (DeliveryArea) TableName() string {
	return "marketplace_delivery_areas"
}
