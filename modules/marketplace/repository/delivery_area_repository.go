package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/domain"
	"gorm.io/gorm"
)

type DeliveryAreaRepository interface {
	Create(row *domain.DeliveryArea) error
	ListByUserID(userID uint) ([]domain.DeliveryArea, error)
	CountByUserID(userID uint) (int64, error)
	GetByIDAndUserID(id, userID uint) (*domain.DeliveryArea, error)
	Update(row *domain.DeliveryArea) error
	Delete(row *domain.DeliveryArea) error
	ResetDefaultByUserID(userID uint) error
	GetRegionByID(id uint) (*adminDomain.Region, error)
	GetDistrictByID(id uint) (*adminDomain.District, error)
	GetMFYByID(id uint) (*adminDomain.MFY, error)
}

type deliveryAreaPostgresRepository struct {
	db *gorm.DB
}

func NewDeliveryAreaRepository(db *gorm.DB) DeliveryAreaRepository {
	return &deliveryAreaPostgresRepository{db: db}
}

func (r *deliveryAreaPostgresRepository) Create(row *domain.DeliveryArea) error {
	return r.db.Create(row).Error
}

func (r *deliveryAreaPostgresRepository) ListByUserID(userID uint) ([]domain.DeliveryArea, error) {
	var rows []domain.DeliveryArea
	err := r.db.Where("user_id = ?", userID).Order("is_default desc, id desc").Find(&rows).Error
	return rows, err
}

func (r *deliveryAreaPostgresRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.DeliveryArea{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *deliveryAreaPostgresRepository) GetByIDAndUserID(id, userID uint) (*domain.DeliveryArea, error) {
	var row domain.DeliveryArea
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *deliveryAreaPostgresRepository) Update(row *domain.DeliveryArea) error {
	return r.db.Save(row).Error
}

func (r *deliveryAreaPostgresRepository) Delete(row *domain.DeliveryArea) error {
	return r.db.Delete(row).Error
}

func (r *deliveryAreaPostgresRepository) ResetDefaultByUserID(userID uint) error {
	return r.db.Model(&domain.DeliveryArea{}).Where("user_id = ?", userID).Update("is_default", false).Error
}

func (r *deliveryAreaPostgresRepository) GetRegionByID(id uint) (*adminDomain.Region, error) {
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

func (r *deliveryAreaPostgresRepository) GetDistrictByID(id uint) (*adminDomain.District, error) {
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

func (r *deliveryAreaPostgresRepository) GetMFYByID(id uint) (*adminDomain.MFY, error) {
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
