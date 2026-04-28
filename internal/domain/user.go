package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	Phone        string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
