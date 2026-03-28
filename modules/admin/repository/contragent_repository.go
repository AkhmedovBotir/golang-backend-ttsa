package repository

import (
	"errors"

	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type ContragentRepository interface {
	Create(row *domain.Contragent) error
	GetPaginated(page, limit int) ([]domain.Contragent, int64, error)
	GetByID(id uint) (*domain.Contragent, error)
	Update(row *domain.Contragent) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
}

type contragentPostgresRepository struct {
	db *gorm.DB
}

func NewContragentRepository(db *gorm.DB) ContragentRepository {
	return &contragentPostgresRepository{db: db}
}

func (r *contragentPostgresRepository) Create(row *domain.Contragent) error {
	return r.db.Create(row).Error
}

func (r *contragentPostgresRepository) GetPaginated(page, limit int) ([]domain.Contragent, int64, error) {
	var rows []domain.Contragent
	var total int64

	if err := r.db.Model(&domain.Contragent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *contragentPostgresRepository) GetByID(id uint) (*domain.Contragent, error) {
	var row domain.Contragent
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentPostgresRepository) Update(row *domain.Contragent) error {
	return r.db.Save(row).Error
}

func (r *contragentPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Contragent{}, id).Error
}

func (r *contragentPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Contragent{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
