package repository

import (
	"errors"

	adminDomain "backend/modules/admin/domain"
	contrDomain "backend/modules/contragents/domain"
	"gorm.io/gorm"
)

type AdminProductRepository interface {
	NextProductCode() (uint64, error)
	Create(row *contrDomain.Product, images []string) error
	GetPaginated(page, limit int, contragentID *uint, moderationStatus *string) ([]contrDomain.Product, int64, error)
	GetByID(id uint) (*contrDomain.Product, error)
	GetImages(productID uint) ([]string, error)
	Update(row *contrDomain.Product, images []string) error
	Delete(id uint) error

	GetCategoryByID(id uint) (*adminDomain.Category, error)
	GetSubcategoryByID(id uint) (*adminDomain.Category, error)
	ContragentExists(id uint) (bool, error)
}

type adminProductPostgresRepository struct {
	db *gorm.DB
}

func NewAdminProductRepository(db *gorm.DB) AdminProductRepository {
	return &adminProductPostgresRepository{db: db}
}

func (r *adminProductPostgresRepository) NextProductCode() (uint64, error) {
	type result struct {
		Max uint64 `gorm:"column:max"`
	}
	var res result
	if err := r.db.Model(&contrDomain.Product{}).Select("COALESCE(MAX(product_code), 0) + 1 AS max").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Max, nil
}

func (r *adminProductPostgresRepository) Create(row *contrDomain.Product, images []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		for i, image := range images {
			img := contrDomain.ProductImage{ProductID: row.ID, Image: image, SortOrder: i}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *adminProductPostgresRepository) GetPaginated(page, limit int, contragentID *uint, moderationStatus *string) ([]contrDomain.Product, int64, error) {
	var rows []contrDomain.Product
	var total int64
	base := r.db.Model(&contrDomain.Product{})
	if contragentID != nil {
		base = base.Where("contragent_id = ?", *contragentID)
	}
	if moderationStatus != nil && *moderationStatus != "" {
		base = base.Where("moderation_status = ?", *moderationStatus)
	}
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

func (r *adminProductPostgresRepository) GetByID(id uint) (*contrDomain.Product, error) {
	var row contrDomain.Product
	err := r.db.First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminProductPostgresRepository) GetImages(productID uint) ([]string, error) {
	var rows []contrDomain.ProductImage
	if err := r.db.Where("product_id = ?", productID).Order("sort_order asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.Image)
	}
	return out, nil
}

func (r *adminProductPostgresRepository) Update(row *contrDomain.Product, images []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		if images != nil {
			if err := tx.Where("product_id = ?", row.ID).Delete(&contrDomain.ProductImage{}).Error; err != nil {
				return err
			}
			for i, image := range images {
				img := contrDomain.ProductImage{ProductID: row.ID, Image: image, SortOrder: i}
				if err := tx.Create(&img).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *adminProductPostgresRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&contrDomain.ProductImage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&contrDomain.Product{}, id).Error
	})
}

func (r *adminProductPostgresRepository) GetCategoryByID(id uint) (*adminDomain.Category, error) {
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

func (r *adminProductPostgresRepository) GetSubcategoryByID(id uint) (*adminDomain.Category, error) {
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

func (r *adminProductPostgresRepository) ContragentExists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&adminDomain.Contragent{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
