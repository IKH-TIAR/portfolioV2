package store

import (
	"context"
	"fmt"
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Store) CreateMessage(ctx context.Context, name, email, body string) (*Message, error) {
	const query = `INSERT INTO messages (name, email, body) VALUES ($1, $2, $3)
	               RETURNING id, created_at`
	var m Message
	err := s.pool.QueryRow(ctx, query, name, email, body).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	m.Name = name
	m.Email = email
	m.Body = body
	return &m, nil
}

func (s *Store) ListMessages(ctx context.Context) ([]Message, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, email, body, created_at FROM messages ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()

	out := make([]Message, 0)
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Body, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) DeleteMessage(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM messages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
