package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user. A duplicate email must surface as
// domain.ErrConflict (Postgres unique_violation, code 23505) so the service
// layer can turn it into a 409 instead of a 500.
func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, name, phone, role, created_at, updated_at)
	values($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Phone,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	var pqErr *pq.Error

	if err != nil {
		switch {
		case errors.As(err, &pqErr):
			if pqErr.Code == "23505" {
				return fmt.Errorf("create user conflict: %w", domain.ErrConflict)
			}
			return err
		default:
			return err
		}
	}
	return nil
}

// GetUserByEmail returns the user with the given email, or domain.ErrNotFound
// when no such user exists.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	query := `SELECT id, email, password_hash, name, phone, role, created_at, updated_at
	FROM users
	WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("get user by email %q: %w", email, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("error get user: %w", err)
	}
	return user, nil
}
