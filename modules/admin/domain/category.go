package domain

import "time"

// Category — parent/subcategory daraxt modeli.
type Category struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `gorm:"size:255;not null" json:"name"`
	Slug       string     `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Image      string     `gorm:"type:text" json:"image"`
	Censored   bool       `gorm:"not null;default:false" json:"censored"`
	ParentID   *uint      `gorm:"index" json:"parent_id,omitempty"`
	Status     string     `gorm:"size:20;not null" json:"status"`
	ExternalID string     `gorm:"size:24;uniqueIndex" json:"external_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
