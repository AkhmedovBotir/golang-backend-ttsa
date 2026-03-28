package repository

import (
	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type ContragentTypeRepository interface {
	Create(row *domain.ContragentType) error
	GetAll() ([]domain.ContragentType, error)
	GetPaginated(page, limit int) ([]domain.ContragentType, int64, error)
	GetByID(id uint) (*domain.ContragentType, error)
	Update(row *domain.ContragentType) error
	Delete(id uint) error
}

type contragentTypePostgresRepository struct {
	db *gorm.DB
}

func NewContragentTypeRepository(db *gorm.DB) ContragentTypeRepository {
	return &contragentTypePostgresRepository{db: db}
}

func (r *contragentTypePostgresRepository) Create(row *domain.ContragentType) error {
	return r.db.Create(row).Error
}

func (r *contragentTypePostgresRepository) GetAll() ([]domain.ContragentType, error) {
	var rows []domain.ContragentType
	return rows, r.db.Order("id asc").Find(&rows).Error
}

func (r *contragentTypePostgresRepository) GetPaginated(page, limit int) ([]domain.ContragentType, int64, error) {
	var rows []domain.ContragentType
	var total int64

	if err := r.db.Model(&domain.ContragentType{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *contragentTypePostgresRepository) GetByID(id uint) (*domain.ContragentType, error) {
	var row domain.ContragentType
	if err := r.db.First(&row, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

func (r *contragentTypePostgresRepository) Update(row *domain.ContragentType) error {
	return r.db.Save(row).Error
}

func (r *contragentTypePostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ContragentType{}, id).Error
}
