package repository

import (
	"errors"

	"backend/modules/admin/domain"
	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(admin *domain.Admin) error
	GetAll() ([]domain.Admin, error)
	GetPaginated(page, limit int) ([]domain.Admin, int64, error)
	GetByID(id uint) (*domain.Admin, error)
	GetByUsername(username string) (*domain.Admin, error)
	Update(admin *domain.Admin) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
	ExistsByUsername(username string, exceptID uint) (bool, error)
}

type adminPostgresRepository struct {
	db *gorm.DB
}

func NewAdminPostgresRepository(db *gorm.DB) AdminRepository {
	return &adminPostgresRepository{db: db}
}

func (r *adminPostgresRepository) Create(admin *domain.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminPostgresRepository) GetAll() ([]domain.Admin, error) {
	var admins []domain.Admin
	err := r.db.Order("id asc").Find(&admins).Error
	return admins, err
}

func (r *adminPostgresRepository) GetPaginated(page, limit int) ([]domain.Admin, int64, error) {
	var admins []domain.Admin
	var total int64

	if err := r.db.Model(&domain.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

func (r *adminPostgresRepository) GetByID(id uint) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.db.First(&admin, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminPostgresRepository) GetByUsername(username string) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminPostgresRepository) Update(admin *domain.Admin) error {
	return r.db.Save(admin).Error
}

func (r *adminPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Admin{}, id).Error
}

func (r *adminPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Admin{}).Where("phone = ?", phone)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *adminPostgresRepository) ExistsByUsername(username string, exceptID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Admin{}).Where("username = ?", username)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
