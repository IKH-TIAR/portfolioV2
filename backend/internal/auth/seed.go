package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type UserStore interface {
	FindUserByEmail(ctx context.Context, email string) (*store.User, error)
	CreateUser(ctx context.Context, email, passwordHash string) (*store.User, error)
}

type Seeder struct {
	store UserStore
}

func NewSeeder(store UserStore) *Seeder {
	return &Seeder{store: store}
}

func (s *Seeder) EnsureAdmin(ctx context.Context, email, password string) error {
	existing, err := s.store.FindUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil
	}
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("look up admin: %w", err)
	}

	hash, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	if _, err := s.store.CreateUser(ctx, email, hash); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}
	return nil
}
