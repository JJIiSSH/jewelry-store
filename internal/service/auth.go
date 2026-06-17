package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo domain.UserRepository
}

const bcryptCost = 12

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register creates a new customer account: it hashes the plaintext password,
// assigns the customer role (a client-provided role is never trusted) and
// persists the user. A duplicate email surfaces as domain.ErrConflict.
func (s *AuthService) Register(ctx context.Context, user domain.User, password string) (uuid.UUID, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return uuid.Nil, fmt.Errorf("password is too long: %w", domain.ErrValidation)
		default:
			return uuid.Nil, err
		}
	}

	now := time.Now()
	user.PasswordHash = string(passHash)
	user.CreatedAt = now
	user.UpdatedAt = now
	user.ID = uuid.New()

	user.Role = domain.RoleCustomer

	err = s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
