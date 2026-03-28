package domain

import "time"

// Manager — menejer (faqat viloyat bilan bog'langan).
type Manager struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:200;not null" json:"name"`
	RegionID             uint      `gorm:"not null;index" json:"viloyat_id"`
	Phone                string    `gorm:"size:13;uniqueIndex;not null" json:"phone"`
	Status               string    `gorm:"size:20;not null" json:"status"`
	Password             string    `gorm:"size:255" json:"-"`
	PasswordSetupAllowed bool      `gorm:"not null;default:true" json:"password_setup_allowed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (Manager) TableName() string {
	return "managers"
}
