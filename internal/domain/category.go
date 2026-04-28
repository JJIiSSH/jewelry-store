package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
