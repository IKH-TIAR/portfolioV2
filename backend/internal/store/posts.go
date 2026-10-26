package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	PostType  string    `json:"postType"`
	Tags      []string  `json:"tags"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostQuery struct {
	PublishedOnly bool
	Tag           string
	Limit         int
	Offset        int
}

func (s *Store) ListPosts(ctx context.Context, q PostQuery) ([]Post, error) {
	query := `SELECT id, slug, title, body, post_type, tags, published, created_at, updated_at
	          FROM posts`
	args := []any{}
	conds := []string{}
	if q.PublishedOnly {
		conds = append(conds, ` published = TRUE`)
	}
	if q.Tag != "" {
		args = append(args, q.Tag)
		conds = append(conds, fmt.Sprintf(` $%d = ANY(tags)`, len(args)))
	}
	if len(conds) > 0 {
		query += ` WHERE` + strings.Join(conds, ` AND`)
	}
	if q.Limit > 0 {
		query += fmt.Sprintf(` ORDER BY created_at DESC, id DESC LIMIT %d OFFSET %d`, q.Limit, q.Offset)
	} else {
		query += ` ORDER BY created_at DESC, id DESC`
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	out := make([]Post, 0)
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Body, &p.PostType, &p.Tags, &p.Published,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetPostBySlug(ctx context.Context, slug string, publishedOnly bool) (*Post, error) {
	query := `SELECT id, slug, title, body, post_type, tags, published, created_at, updated_at
	          FROM posts WHERE slug = $1`
	if publishedOnly {
		query += ` AND published = TRUE`
	}
	return scanPost(s.pool.QueryRow(ctx, query, slug))
}

func (s *Store) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	const query = `SELECT id, slug, title, body, post_type, tags, published, created_at, updated_at
	               FROM posts WHERE id = $1`
	return scanPost(s.pool.QueryRow(ctx, query, id))
}

func scanPost(row pgx.Row) (*Post, error) {
	var p Post
	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Body, &p.PostType, &p.Tags, &p.Published,
		&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan post: %w", err)
	}
	return &p, nil
}

func (s *Store) CreatePost(ctx context.Context, p Post) (*Post, error) {
	const query = `INSERT INTO posts (slug, title, body, post_type, tags, published)
	               VALUES ($1, $2, $3, $4, $5, $6)
	               RETURNING id, created_at, updated_at`
	err := s.pool.QueryRow(ctx, query, p.Slug, p.Title, p.Body, p.PostType, p.Tags, p.Published).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return &p, nil
}

func (s *Store) UpdatePost(ctx context.Context, id int64, p Post) (*Post, error) {
	const query = `UPDATE posts SET slug = $2, title = $3, body = $4, post_type = $5,
	                      tags = $6, published = $7, updated_at = now()
	               WHERE id = $1
	               RETURNING created_at, updated_at`
	err := s.pool.QueryRow(ctx, query, id, p.Slug, p.Title, p.Body, p.PostType, p.Tags, p.Published).
		Scan(&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	p.ID = id
	return &p, nil
}

func (s *Store) DeletePost(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
