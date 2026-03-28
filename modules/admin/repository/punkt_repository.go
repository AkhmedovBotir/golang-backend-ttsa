package repository

import (
	"errors"

	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type PunktRepository interface {
	Create(row *domain.Punkt) error
	GetPaginated(page, limit int) ([]domain.Punkt, int64, error)
	GetByID(id uint) (*domain.Punkt, error)
	Update(row *domain.Punkt) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
}

type punktPostgresRepository struct {
	db *gorm.DB
}

func NewPunktRepository(db *gorm.DB) PunktRepository {
	return &punktPostgresRepository{db: db}
}

func (r *punktPostgresRepository) Create(row *domain.Punkt) error {
	return r.db.Create(row).Error
}

func (r *punktPostgresRepository) GetPaginated(page, limit int) ([]domain.Punkt, int64, error) {
	var rows []domain.Punkt
	var total int64

	if err := r.db.Model(&domain.Punkt{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *punktPostgresRepository) GetByID(id uint) (*domain.Punkt, error) {
	var row domain.Punkt
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *punktPostgresRepository) Update(row *domain.Punkt) error {
	return r.db.Save(row).Error
}

func (r *punktPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Punkt{}, id).Error
}

func (r *punktPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Punkt{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
