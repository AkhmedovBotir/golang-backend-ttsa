package repository

import (
	"errors"

	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type NeighborhoodShopRepository interface {
	Create(row *domain.NeighborhoodShop) error
	GetPaginated(page, limit int) ([]domain.NeighborhoodShop, int64, error)
	GetByID(id uint) (*domain.NeighborhoodShop, error)
	Update(row *domain.NeighborhoodShop) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
}

type neighborhoodShopPostgresRepository struct {
	db *gorm.DB
}

func NewNeighborhoodShopRepository(db *gorm.DB) NeighborhoodShopRepository {
	return &neighborhoodShopPostgresRepository{db: db}
}

func (r *neighborhoodShopPostgresRepository) Create(row *domain.NeighborhoodShop) error {
	return r.db.Create(row).Error
}

func (r *neighborhoodShopPostgresRepository) GetPaginated(page, limit int) ([]domain.NeighborhoodShop, int64, error) {
	var rows []domain.NeighborhoodShop
	var total int64

	if err := r.db.Model(&domain.NeighborhoodShop{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *neighborhoodShopPostgresRepository) GetByID(id uint) (*domain.NeighborhoodShop, error) {
	var row domain.NeighborhoodShop
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *neighborhoodShopPostgresRepository) Update(row *domain.NeighborhoodShop) error {
	return r.db.Save(row).Error
}

func (r *neighborhoodShopPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.NeighborhoodShop{}, id).Error
}

func (r *neighborhoodShopPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.NeighborhoodShop{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
