package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string
type ProductType string

const (
	ProductTypeRing     ProductType = "ring"
	ProductTypeNecklace ProductType = "necklace"
	ProductTypeBracelet ProductType = "bracelet"
	ProductTypeEarrings ProductType = "earrings"
	ProductTypeBrooch   ProductType = "brooch"
	ProductTypePendant  ProductType = "pendant"

	ProductStatusPublished ProductStatus = "published"
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusSold      ProductStatus = "sold"
	ProductStatusArchived  ProductStatus = "archived"
)

type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Title       string
	Slug        string
	Description string
	Story       string
	Stone       string
	Materials   []string
	Price       float64
	WeightG     float64
	Size        string
	IsUnique    bool
	Stock       int
	Status      ProductStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
