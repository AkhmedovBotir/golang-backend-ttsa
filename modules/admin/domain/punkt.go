package domain

import "time"

// Punkt — punkt (viloyat va tuman bilan bog'langan).
type Punkt struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:200;not null" json:"name"`
	RegionID             uint      `gorm:"not null;index" json:"viloyat_id"`
	DistrictID           uint      `gorm:"not null;index" json:"tuman_id"`
	Phone                string    `gorm:"size:13;uniqueIndex;not null" json:"phone"`
	Status               string    `gorm:"size:20;not null" json:"status"`
	Password             string    `gorm:"size:255" json:"-"`
	PasswordSetupAllowed bool      `gorm:"not null;default:true" json:"password_setup_allowed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (Punkt) TableName() string {
	return "punkts"
}
