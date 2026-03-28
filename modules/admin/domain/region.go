package domain

import "time"

type Region struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ExternalID string    `gorm:"size:24;uniqueIndex" json:"external_id"`
	Name       string    `gorm:"size:160;not null" json:"name"`
	Code       string    `gorm:"size:160;not null" json:"code"`
	Status     string    `gorm:"size:20;not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type District struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ExternalID string    `gorm:"size:24;uniqueIndex" json:"external_id"`
	RegionID   uint      `gorm:"index;not null" json:"region_id"`
	Name       string    `gorm:"size:160;not null" json:"name"`
	Code       string    `gorm:"size:160;not null" json:"code"`
	Status     string    `gorm:"size:20;not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MFY struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ExternalID string    `gorm:"size:24;uniqueIndex" json:"external_id"`
	DistrictID uint      `gorm:"index;not null" json:"district_id"`
	Name       string    `gorm:"size:180;not null" json:"name"`
	Code       string    `gorm:"size:180;not null" json:"code"`
	Status     string    `gorm:"size:20;not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
