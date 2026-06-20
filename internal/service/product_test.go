package service

import (
	"context"
	"errors"
	"testing"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockProductRepo struct {
	returnerErr  error
	capturedItem domain.Product
}

func (m *mockProductRepo) CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error) {

	m.capturedItem = item
	return uuid.New(), m.returnerErr
}

func (m *mockProductRepo) GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var pr domain.Product
	return pr, m.returnerErr
}

func (m *mockProductRepo) GetProducts(ctx context.Context, queryParams domain.ProductFilter) (domain.ProductList, error) {
	var prs domain.ProductList
	return prs, m.returnerErr
}

func (m *mockProductRepo) UpdateProduct(ctx context.Context, item domain.Product) error {
	return m.returnerErr
}

func (m *mockProductRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return m.returnerErr
}

func (m *mockProductRepo) ChangeProductStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error {
	return m.returnerErr
}

func TestProductService_CreateProduct_SetsSlugAndTimestamps(t *testing.T) {

	var item domain.Product

	item.Title = "Amethyst Ring"

	mock := &mockProductRepo{}
	s := NewProductService(mock)

	id, err := s.CreateProduct(context.Background(), item)

	assert.Equal(t, "amethyst-ring", mock.capturedItem.Slug)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	assert.False(t, mock.capturedItem.CreatedAt.IsZero())
	assert.False(t, mock.capturedItem.UpdatedAt.IsZero())

}

func TestProductService_CreateProduct_ReturnsErrorOnRepoFailure(t *testing.T) {

	var item domain.Product

	item.Title = "test"

	mock := &mockProductRepo{returnerErr: errors.New("not created")}

	s := NewProductService(mock)

	_, err := s.CreateProduct(context.Background(), item)

	assert.Error(t, err)

}
