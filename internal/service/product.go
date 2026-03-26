package service

import (
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"github.com/jonathanCaamano/portfolio-gateway/internal/repository"
)

type ProductService struct {
	repo  *repository.ProductRepo
	minio *MinioService
}

func NewProductService(repo *repository.ProductRepo, minio *MinioService) *ProductService {
	return &ProductService{repo, minio}
}

func (s *ProductService) List(f repository.ProductFilter) ([]domain.Product, int64, error) {
	return s.repo.List(f)
}

func (s *ProductService) Get(id uuid.UUID) (*domain.Product, error) { return s.repo.FindByID(id) }
func (s *ProductService) Create(p *domain.Product) error            { return s.repo.Create(p) }
func (s *ProductService) Update(p *domain.Product) error            { return s.repo.Update(p) }
func (s *ProductService) Delete(id uuid.UUID) error                  { return s.repo.Delete(id) }
func (s *ProductService) Stats() (map[string]any, error)             { return s.repo.Stats() }

func (s *ProductService) UploadImage(productID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	contentType := header.Header.Get("Content-Type")
	return s.minio.Upload(productID.String()+"/"+header.Filename, data, contentType)
}
