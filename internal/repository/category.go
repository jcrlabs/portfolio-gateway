package repository

import (
	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db} }

func (r *CategoryRepo) List() ([]domain.Category, error) {
	var cats []domain.Category
	return cats, r.db.Find(&cats).Error
}

func (r *CategoryRepo) FindByID(id uuid.UUID) (*domain.Category, error) {
	var c domain.Category
	return &c, r.db.First(&c, "id = ?", id).Error
}

func (r *CategoryRepo) Create(c *domain.Category) error {
	return r.db.Create(c).Error
}

func (r *CategoryRepo) Update(c *domain.Category) error {
	return r.db.Save(c).Error
}

func (r *CategoryRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Category{}, "id = ?", id).Error
}
