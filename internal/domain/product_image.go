package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	URL        string
	Alt        string
	OrderIndex int
	IsPrimary  bool
	CreatedAt  time.Time
}
