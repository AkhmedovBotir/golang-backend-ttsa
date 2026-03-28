package repository

import (
	"errors"

	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type ManagerRepository interface {
	Create(row *domain.Manager) error
	GetPaginated(page, limit int) ([]domain.Manager, int64, error)
	GetByID(id uint) (*domain.Manager, error)
	Update(row *domain.Manager) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
}

type managerPostgresRepository struct {
	db *gorm.DB
}

func NewManagerRepository(db *gorm.DB) ManagerRepository {
	return &managerPostgresRepository{db: db}
}

func (r *managerPostgresRepository) Create(row *domain.Manager) error {
	return r.db.Create(row).Error
}

func (r *managerPostgresRepository) GetPaginated(page, limit int) ([]domain.Manager, int64, error) {
	var rows []domain.Manager
	var total int64

	if err := r.db.Model(&domain.Manager{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *managerPostgresRepository) GetByID(id uint) (*domain.Manager, error) {
	var row domain.Manager
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *managerPostgresRepository) Update(row *domain.Manager) error {
	return r.db.Save(row).Error
}

func (r *managerPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Manager{}, id).Error
}

func (r *managerPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Manager{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
