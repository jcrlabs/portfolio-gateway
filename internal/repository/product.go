package repository

import (
	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"gorm.io/gorm"
)

type ProductRepo struct{ db *gorm.DB }

func NewProductRepo(db *gorm.DB) *ProductRepo { return &ProductRepo{db} }

type ProductFilter struct {
	Search     string
	CategoryID *uuid.UUID
	Page       int
	Limit      int
}

func (r *ProductRepo) List(f ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	q := r.db.Model(&domain.Product{}).Preload("Category")
	if f.Search != "" {
		q = q.Where("name ILIKE ?", "%"+f.Search+"%")
	}
	if f.CategoryID != nil {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	q.Count(&total)

	page, limit := f.Page, f.Limit
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return products, total, q.Offset((page-1)*limit).Limit(limit).Find(&products).Error
}

func (r *ProductRepo) FindByID(id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	return &p, r.db.Preload("Category").First(&p, "id = ?", id).Error
}

func (r *ProductRepo) Create(p *domain.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepo) Update(p *domain.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Product{}, "id = ?", id).Error
}

func (r *ProductRepo) Stats() (map[string]any, error) {
	var totalProducts int64
	var stockValue float64
	var lowStock int64

	r.db.Model(&domain.Product{}).Count(&totalProducts)
	r.db.Model(&domain.Product{}).Select("COALESCE(SUM(price * stock), 0)").Scan(&stockValue)
	r.db.Model(&domain.Product{}).Where("stock < 10").Count(&lowStock)

	return map[string]any{
		"total_products": totalProducts,
		"stock_value":    stockValue,
		"low_stock_alerts": lowStock,
	}, nil
}
