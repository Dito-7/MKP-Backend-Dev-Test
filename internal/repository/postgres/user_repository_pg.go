package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
)

type userRepoPG struct {
	db *sql.DB
}

func NewUserRepositoryPG(db *sql.DB) repository.UserRepository {
	return &userRepoPG{db: db}
}

func (r *userRepoPG) Create(ctx context.Context, u *entity.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, phone_number, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		u.FullName, u.Email, u.PasswordHash, u.PhoneNumber, u.Role, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *userRepoPG) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, phone_number, role, is_active, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1) AND is_active = true
	`
	var u entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.PhoneNumber, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	return &u, nil
}

func (r *userRepoPG) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, phone_number, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`
	var u entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.PhoneNumber, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}
	return &u, nil
}
