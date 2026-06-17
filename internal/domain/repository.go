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

// ProductList is a page of products plus Total — the number of items matching
// the filter ignoring LIMIT/OFFSET. It backs the paginated list response.
type ProductList struct {
	Items []Product
	Total int
}

type ProductRepository interface {
	GetProducts(ctx context.Context, queryParams ProductFilter) (ProductList, error)
	GetProductByID(ctx context.Context, ID uuid.UUID) (Product, error)
	CreateProduct(ctx context.Context, item Product) (uuid.UUID, error)
	UpdateProduct(ctx context.Context, item Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ChangeProductStatus(ctx context.Context, id uuid.UUID, status ProductStatus) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
