package domain

import "time"

type ContragentType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ExternalID string   `gorm:"size:24;uniqueIndex" json:"external_id"`
	Name      string    `gorm:"size:160;not null" json:"name"`
	Icon      string    `gorm:"size:255;not null" json:"icon"`
	Status    string    `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
