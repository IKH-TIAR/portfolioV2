package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	const q = `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	var u User
	err := s.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	const q = `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at`

	var u User
	err := s.pool.QueryRow(ctx, q, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}
