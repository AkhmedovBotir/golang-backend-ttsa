package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type ContragentProductRepository interface {
	NextProductCode() (uint64, error)
	Create(row *domain.Product, images []string) error
	GetPaginatedByContragent(contragentID uint, page, limit int) ([]domain.Product, int64, error)
	GetByIDAndContragent(id, contragentID uint) (*domain.Product, error)
	GetImages(productID uint) ([]string, error)
	Update(row *domain.Product, images []string) error
	Delete(id uint) error

	GetCategoryByID(id uint) (*adminDomain.Category, error)
	GetSubcategoryByID(id uint) (*adminDomain.Category, error)
}

type contragentProductPostgresRepository struct {
	db *gorm.DB
}

func NewContragentProductRepository(db *gorm.DB) ContragentProductRepository {
	return &contragentProductPostgresRepository{db: db}
}

func (r *contragentProductPostgresRepository) NextProductCode() (uint64, error) {
	type result struct {
		Max uint64 `gorm:"column:max"`
	}
	var res result
	if err := r.db.Model(&domain.Product{}).Select("COALESCE(MAX(product_code), 0) + 1 AS max").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Max, nil
}

func (r *contragentProductPostgresRepository) Create(row *domain.Product, images []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		for i, image := range images {
			img := domain.ProductImage{ProductID: row.ID, Image: image, SortOrder: i}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *contragentProductPostgresRepository) GetPaginatedByContragent(contragentID uint, page, limit int) ([]domain.Product, int64, error) {
	var rows []domain.Product
	var total int64
	base := r.db.Model(&domain.Product{}).Where("contragent_id = ?", contragentID)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := base.Order("id desc").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *contragentProductPostgresRepository) GetByIDAndContragent(id, contragentID uint) (*domain.Product, error) {
	var row domain.Product
	err := r.db.Where("id = ? AND contragent_id = ?", id, contragentID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentProductPostgresRepository) GetImages(productID uint) ([]string, error) {
	var rows []domain.ProductImage
	if err := r.db.Where("product_id = ?", productID).Order("sort_order asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.Image)
	}
	return out, nil
}

func (r *contragentProductPostgresRepository) Update(row *domain.Product, images []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		if images != nil {
			if err := tx.Where("product_id = ?", row.ID).Delete(&domain.ProductImage{}).Error; err != nil {
				return err
			}
			for i, image := range images {
				img := domain.ProductImage{ProductID: row.ID, Image: image, SortOrder: i}
				if err := tx.Create(&img).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *contragentProductPostgresRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&domain.ProductImage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Product{}, id).Error
	})
}

func (r *contragentProductPostgresRepository) GetCategoryByID(id uint) (*adminDomain.Category, error) {
	var row adminDomain.Category
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contragentProductPostgresRepository) GetSubcategoryByID(id uint) (*adminDomain.Category, error) {
	var row adminDomain.Category
	err := r.db.Where("id = ? AND parent_id IS NOT NULL", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
