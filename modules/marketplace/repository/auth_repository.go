package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByPhone(phone string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	CreateUser(row *domain.User) error
	UpdateUser(row *domain.User) error
	CreateCode(row *domain.VerificationCode) error
	GetLatestActiveCode(phone, purpose string) (*domain.VerificationCode, error)
	UpdateCode(row *domain.VerificationCode) error
	GetRegionByID(id uint) (*adminDomain.Region, error)
	GetDistrictByID(id uint) (*adminDomain.District, error)
	GetMFYByID(id uint) (*adminDomain.MFY, error)
}

type authPostgresRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authPostgresRepository{db: db}
}

func (r *authPostgresRepository) GetUserByPhone(phone string) (*domain.User, error) {
	var row domain.User
	err := r.db.Where("phone = ?", phone).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *authPostgresRepository) GetUserByID(id uint) (*domain.User, error) {
	var row domain.User
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *authPostgresRepository) CreateUser(row *domain.User) error {
	return r.db.Create(row).Error
}

func (r *authPostgresRepository) UpdateUser(row *domain.User) error {
	return r.db.Save(row).Error
}

func (r *authPostgresRepository) CreateCode(row *domain.VerificationCode) error {
	return r.db.Create(row).Error
}

func (r *authPostgresRepository) GetLatestActiveCode(phone, purpose string) (*domain.VerificationCode, error) {
	var row domain.VerificationCode
	err := r.db.Where("phone = ? AND purpose = ? AND used_at IS NULL", phone, purpose).Order("id desc").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *authPostgresRepository) UpdateCode(row *domain.VerificationCode) error {
	return r.db.Save(row).Error
}

func (r *authPostgresRepository) GetRegionByID(id uint) (*adminDomain.Region, error) {
	var row adminDomain.Region
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *authPostgresRepository) GetDistrictByID(id uint) (*adminDomain.District, error) {
	var row adminDomain.District
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *authPostgresRepository) GetMFYByID(id uint) (*adminDomain.MFY, error) {
	var row adminDomain.MFY
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
