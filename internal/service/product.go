package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error) {

	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.Slug = slug.Make(item.Title)

	if item.Status == "" {
		item.Status = domain.ProductStatusDraft
	}

	_, err := s.repo.CreateProduct(ctx, item)

	if err != nil {
		return item.ID, fmt.Errorf("failed to create product: %w", err)
	}
	return item.ID, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) GetProducts(ctx context.Context, queryParams domain.ProductFilter) ([]domain.Product, error) {

	return s.repo.GetProducts(ctx, queryParams)
}

func (s *ProductService) UpdateProduct(ctx context.Context, item domain.Product) error {
	item.Slug = slug.Make(item.Title)
	item.UpdatedAt = time.Now()
	return s.repo.UpdateProduct(ctx, item)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) ChangeProductStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error {
	return s.repo.ChangeProductStatus(ctx, id, status)
}
