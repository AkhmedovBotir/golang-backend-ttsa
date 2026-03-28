package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type ContragentAuthRepository interface {
	GetContragentByPhone(phone string) (*adminDomain.Contragent, error)
	GetContragentByID(id uint) (*adminDomain.Contragent, error)
	UpdateContragent(row *adminDomain.Contragent) error
	CreateCode(row *domain.VerificationCode) error
	GetLatestActiveCode(phone, purpose string) (*domain.VerificationCode, error)
	UpdateCode(row *domain.VerificationCode) error
}

type contragentAuthPostgresRepository struct {
	db *gorm.DB
}

func NewContragentAuthRepository(db *gorm.DB) ContragentAuthRepository {
	return &contragentAuthPostgresRepository{db: db}
}

func (r *contragentAuthPostgresRepository) GetContragentByPhone(phone string) (*adminDomain.Contragent, error) {
	var row adminDomain.Contragent
	err := r.db.Where("phone = ?", phone).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentAuthPostgresRepository) GetContragentByID(id uint) (*adminDomain.Contragent, error) {
	var row adminDomain.Contragent
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentAuthPostgresRepository) UpdateContragent(row *adminDomain.Contragent) error {
	return r.db.Save(row).Error
}

func (r *contragentAuthPostgresRepository) CreateCode(row *domain.VerificationCode) error {
	return r.db.Create(row).Error
}

func (r *contragentAuthPostgresRepository) GetLatestActiveCode(phone, purpose string) (*domain.VerificationCode, error) {
	var row domain.VerificationCode
	err := r.db.
		Where("phone = ? AND purpose = ? AND used_at IS NULL", phone, purpose).
		Order("id desc").
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentAuthPostgresRepository) UpdateCode(row *domain.VerificationCode) error {
	return r.db.Save(row).Error
}
