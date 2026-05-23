package dto

import (
	"time"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Title       string               `json:"title" binding:"required,min=3,max=200"`
	Price       float64              `json:"price" binding:"required,gt=0"`
	Description string               `json:"description" binding:"min=10,max=1000"`
	Story       string               `json:"story" binding:"min=10,max=1000"`
	Stone       string               `json:"stone" binding:"required,min=3,max=100"`
	Materials   []string             `json:"materials" binding:"required,min=1"`
	WeightG     float64              `json:"weight_g" binding:"required"`
	Size        string               `json:"size" binding:"required,min=2,max=10"`
	IsUnique    bool                 `json:"is_unique"`
	Stock       int                  `json:"stock" binding:"gte=1"`
	Status      domain.ProductStatus `json:"status"`
	CategoryID  uuid.UUID            `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Description string    `json:"description" binding:"min=10,max=1000"`
	Story       string    `json:"story" binding:"min=10,max=1000"`
	Stone       string    `json:"stone" binding:"required,min=3,max=100"`
	Materials   []string  `json:"materials" binding:"required,min=1"`
	WeightG     float64   `json:"weight_g" binding:"required"`
	Size        string    `json:"size" binding:"required,min=2,max=10"`
	IsUnique    bool      `json:"is_unique"`
	Stock       int       `json:"stock" binding:"gte=1"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
}

type ChangeStatusRequest struct {
	Status domain.ProductStatus `json:"status" binding:"required,oneof=draft published sold archived"`
}

type ListProductsQuery struct {
	Category string  `form:"category"`
	Stone    string  `form:"stone"`
	MinPrice float64 `form:"min_price" binding:"gte=0"`
	MaxPrice float64 `form:"max_price" binding:"gte=0"`
	Sort     string  `form:"sort" binding:"omitempty,oneof=price_asc price_desc newest"`
	Page     int     `form:"page,default=1" binding:"gte=1"`
	Limit    int     `form:"limit,default=20" binding:"gte=1"`
}

type ProductResponse struct {
	ID          uuid.UUID            `json:"id"`
	CategoryID  uuid.UUID            `json:"category_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Story       string               `json:"story"`
	Stone       string               `json:"stone"`
	Materials   []string             `json:"materials"`
	Price       float64              `json:"price"`
	WeightG     float64              `json:"weight_g"`
	Size        string               `json:"size"`
	IsUnique    bool                 `json:"is_unique"`
	Stock       int                  `json:"stock"`
	Status      domain.ProductStatus `json:"status"`
	Slug        string               `json:"slug"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}
