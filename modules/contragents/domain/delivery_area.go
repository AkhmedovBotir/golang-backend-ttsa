package domain

import "time"

type ContragentDeliveryRegion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ContragentID uint      `gorm:"index:idx_contragent_region,unique;not null" json:"contragent_id"`
	RegionID     uint      `gorm:"index:idx_contragent_region,unique;not null" json:"region_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ContragentDeliveryRegion) TableName() string {
	return "contragent_delivery_regions"
}

type ContragentDeliveryDistrict struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ContragentID uint      `gorm:"index:idx_contragent_district,unique;not null" json:"contragent_id"`
	DistrictID   uint      `gorm:"index:idx_contragent_district,unique;not null" json:"district_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ContragentDeliveryDistrict) TableName() string {
	return "contragent_delivery_districts"
}
