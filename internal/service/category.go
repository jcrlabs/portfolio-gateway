package service

import (
	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"github.com/jonathanCaamano/portfolio-gateway/internal/repository"
)

type CategoryService struct{ repo *repository.CategoryRepo }

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) List() ([]domain.Category, error)          { return s.repo.List() }
func (s *CategoryService) Get(id uuid.UUID) (*domain.Category, error) { return s.repo.FindByID(id) }
func (s *CategoryService) Create(c *domain.Category) error            { return s.repo.Create(c) }
func (s *CategoryService) Update(c *domain.Category) error            { return s.repo.Update(c) }
func (s *CategoryService) Delete(id uuid.UUID) error                  { return s.repo.Delete(id) }
