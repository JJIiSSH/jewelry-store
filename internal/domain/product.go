package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusPublished ProductStatus = "published"
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusSold      ProductStatus = "sold"
	ProductStatusArchived  ProductStatus = "archived"
)

type Product struct {
	ID          uuid.UUID     `db:"id"`
	CategoryID  uuid.UUID     `db:"category_id"`
	Title       string        `db:"title"`
	Slug        string        `db:"slug"`
	Description string        `db:"description"`
	Story       string        `db:"story"`
	Stone       string        `db:"stone"`
	Materials   []string      `db:"materials"`
	Price       float64       `db:"price"`
	WeightG     float64       `db:"weight_g"`
	Size        string        `db:"size"`
	IsUnique    bool          `db:"is_unique"`
	Stock       int           `db:"stock"`
	Status      ProductStatus `db:"status"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}
