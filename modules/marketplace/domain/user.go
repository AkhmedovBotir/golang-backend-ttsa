package domain

import "time"

const (
	GenderMale   = "erkak"
	GenderFemale = "ayol"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Phone      string    `gorm:"size:13;uniqueIndex;not null" json:"phone"`
	FirstName  string    `gorm:"size:120;not null" json:"first_name"`
	LastName   string    `gorm:"size:120;not null" json:"last_name"`
	Gender     string    `gorm:"size:16;not null" json:"gender"`
	Avatar     string    `gorm:"type:text" json:"avatar,omitempty"`
	RegionID   uint      `gorm:"not null;index" json:"region_id"`
	DistrictID uint      `gorm:"not null;index" json:"district_id"`
	MFYID      uint      `gorm:"column:mfy_id;not null;index" json:"mfy_id"`
	BirthDate  time.Time `gorm:"not null" json:"birth_date"`
	Status     string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "marketplace_users"
}
