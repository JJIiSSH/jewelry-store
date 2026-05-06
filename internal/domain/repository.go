package domain

import (
	"context"

	"github.com/google/uuid"
)

type Sort string

const (
	SortPriceASC  Sort = "price_asc"
	SortPriceDECS Sort = "price_desc"
	SortNewest    Sort = "newest"
)

type ProductFilter struct {
	Category string
	Stone    string
	MinPrice float64
	MaxPrice float64
	Sort     Sort
	Page     int
	Limit    int
}

type ProductRepository interface {
	GetProducts(ctx context.Context, queryParams ProductFilter) ([]Product, error)
	GetProductByID(ctx context.Context, ID uuid.UUID) (Product, error)
	CreateProduct(ctx context.Context, item Product) (uuid.UUID, error)
	UpdateProduct(ctx context.Context, item Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ChangeProductStatus(ctx context.Context, id uuid.UUID, status ProductStatus) error
}
