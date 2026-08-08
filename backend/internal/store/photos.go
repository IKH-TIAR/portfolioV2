package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Photo struct {
	ID        int64      `json:"id"`
	URL       string     `json:"url"`
	Caption   string     `json:"caption"`
	Location  string     `json:"location"`
	TakenAt   *time.Time `json:"takenAt,omitempty"`
	Published bool       `json:"published"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (s *Store) ListPhotos(ctx context.Context, publishedOnly bool) ([]Photo, error) {
	query := `SELECT id, url, caption, location, taken_at, published, created_at
	          FROM photos`
	args := []any{}
	if publishedOnly {
		query += ` WHERE published = TRUE`
	}
	query += ` ORDER BY taken_at DESC NULLS LAST, created_at DESC, id DESC`

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}
	defer rows.Close()

	out := make([]Photo, 0)
	for rows.Next() {
		var p Photo
		if err := rows.Scan(&p.ID, &p.URL, &p.Caption, &p.Location, &p.TakenAt, &p.Published,
			&p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan photo: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetPhotoByID(ctx context.Context, id int64) (*Photo, error) {
	const query = `SELECT id, url, caption, location, taken_at, published, created_at
	               FROM photos WHERE id = $1`
	var p Photo
	err := s.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.URL, &p.Caption, &p.Location, &p.TakenAt,
		&p.Published, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan photo: %w", err)
	}
	return &p, nil
}

func (s *Store) CreatePhoto(ctx context.Context, p Photo) (*Photo, error) {
	const query = `INSERT INTO photos (url, caption, location, taken_at, published)
	               VALUES ($1, $2, $3, $4, $5)
	               RETURNING id, created_at`
	err := s.pool.QueryRow(ctx, query, p.URL, p.Caption, p.Location, p.TakenAt, p.Published).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create photo: %w", err)
	}
	return &p, nil
}

func (s *Store) UpdatePhoto(ctx context.Context, id int64, p Photo) (*Photo, error) {
	const query = `UPDATE photos SET url = $2, caption = $3, location = $4, taken_at = $5,
	                      published = $6
	               WHERE id = $1
	               RETURNING created_at`
	err := s.pool.QueryRow(ctx, query, id, p.URL, p.Caption, p.Location, p.TakenAt, p.Published).
		Scan(&p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update photo: %w", err)
	}
	p.ID = id
	return &p, nil
}

func (s *Store) DeletePhoto(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete photo: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
