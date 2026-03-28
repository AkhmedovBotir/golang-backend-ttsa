package domain

import "time"

type VerificationCode struct {
	ID           uint       `gorm:"primaryKey"`
	ContragentID uint       `gorm:"not null;index"`
	Phone        string     `gorm:"size:13;not null;index"`
	Code         string     `gorm:"size:10;not null"`
	Purpose      string     `gorm:"size:50;not null;index"`
	ExpiresAt    time.Time  `gorm:"not null;index"`
	VerifiedAt   *time.Time `gorm:"index"`
	UsedAt       *time.Time `gorm:"index"`
	Attempts     int        `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (VerificationCode) TableName() string {
	return "contragent_verification_codes"
}
