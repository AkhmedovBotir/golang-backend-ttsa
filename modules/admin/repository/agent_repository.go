package repository

import (
	"errors"

	"backend/modules/admin/domain"

	"gorm.io/gorm"
)

type AgentRepository interface {
	Create(row *domain.Agent) error
	GetPaginated(page, limit int) ([]domain.Agent, int64, error)
	GetByID(id uint) (*domain.Agent, error)
	Update(row *domain.Agent) error
	Delete(id uint) error
	ExistsByPhone(phone string, exceptID uint) (bool, error)
}

type agentPostgresRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentPostgresRepository{db: db}
}

func (r *agentPostgresRepository) Create(row *domain.Agent) error {
	return r.db.Create(row).Error
}

func (r *agentPostgresRepository) GetPaginated(page, limit int) ([]domain.Agent, int64, error) {
	var rows []domain.Agent
	var total int64

	if err := r.db.Model(&domain.Agent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("id asc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *agentPostgresRepository) GetByID(id uint) (*domain.Agent, error) {
	var row domain.Agent
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *agentPostgresRepository) Update(row *domain.Agent) error {
	return r.db.Save(row).Error
}

func (r *agentPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Agent{}, id).Error
}

func (r *agentPostgresRepository) ExistsByPhone(phone string, exceptID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Agent{}).Where("phone = ?", phone)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
