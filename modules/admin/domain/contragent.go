package domain

import "time"

// Contragent — tashkilot (viloyat, tuman, MFY va faoliyat turi bilan bog'langan).
type Contragent struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:200;not null" json:"name"`
	INN                  string    `gorm:"size:12;not null;index" json:"inn"`
	RegionID             uint      `gorm:"not null;index" json:"region_id"`
	DistrictID           uint      `gorm:"not null;index" json:"district_id"`
	MFYID                uint      `gorm:"column:mfy_id;not null;index" json:"mfy_id"`
	Phone                string    `gorm:"size:13;uniqueIndex;not null" json:"phone"`
	Logo                 string    `gorm:"type:text" json:"logo,omitempty"`
	ActivityTypeID       uint      `gorm:"column:activity_type_id;not null;index" json:"activity_type_id"`
	Status               string    `gorm:"size:20;not null" json:"status"`
	Password             string    `gorm:"size:255" json:"-"`
	PasswordSetupAllowed bool      `gorm:"not null;default:true" json:"password_setup_allowed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TableName — jadval nomi (oldingi `district_contragents` dan farqli; kerak bo'lsa ma'lumotni migratsiya qiling).
func (Contragent) TableName() string {
	return "contragents"
}
