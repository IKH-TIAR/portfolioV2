package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Project struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Tagline     string    `json:"tagline"`
	Description string    `json:"description"`
	Tech        []string  `json:"tech"`
	RepoURL     string    `json:"repoUrl"`
	DemoURL     string    `json:"demoUrl"`
	ImageURL    string    `json:"imageUrl"`
	Featured    bool      `json:"featured"`
	SortOrder   int       `json:"sortOrder"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectQuery struct {
	PublishedOnly bool
	Limit         int
	Offset        int
}

func (s *Store) ListProjects(ctx context.Context, q ProjectQuery) ([]Project, error) {
	query := `SELECT id, slug, title, tagline, description, tech, repo_url, demo_url, image_url,
	                 featured, sort_order, published, created_at, updated_at
	          FROM projects`
	args := []any{}
	if q.PublishedOnly {
		query += ` WHERE published = TRUE`
	}
	if q.Limit > 0 {
		query += fmt.Sprintf(` ORDER BY sort_order ASC, id ASC LIMIT %d OFFSET %d`, q.Limit, q.Offset)
	} else {
		query += ` ORDER BY sort_order ASC, id ASC`
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	out := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Tagline, &p.Description, &p.Tech,
			&p.RepoURL, &p.DemoURL, &p.ImageURL, &p.Featured, &p.SortOrder, &p.Published,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetProjectBySlug(ctx context.Context, slug string, publishedOnly bool) (*Project, error) {
	query := `SELECT id, slug, title, tagline, description, tech, repo_url, demo_url, image_url,
	                 featured, sort_order, published, created_at, updated_at
	          FROM projects WHERE slug = $1`
	if publishedOnly {
		query += ` AND published = TRUE`
	}
	return scanProject(s.pool.QueryRow(ctx, query, slug))
}

func (s *Store) GetProjectByID(ctx context.Context, id int64) (*Project, error) {
	const query = `SELECT id, slug, title, tagline, description, tech, repo_url, demo_url, image_url,
	                      featured, sort_order, published, created_at, updated_at
	               FROM projects WHERE id = $1`
	return scanProject(s.pool.QueryRow(ctx, query, id))
}

func scanProject(row pgx.Row) (*Project, error) {
	var p Project
	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Tagline, &p.Description, &p.Tech,
		&p.RepoURL, &p.DemoURL, &p.ImageURL, &p.Featured, &p.SortOrder, &p.Published,
		&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan project: %w", err)
	}
	return &p, nil
}

func (s *Store) CreateProject(ctx context.Context, p Project) (*Project, error) {
	const query = `INSERT INTO projects (slug, title, tagline, description, tech, repo_url, demo_url,
	                                    image_url, featured, sort_order, published)
	               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	               RETURNING id, created_at, updated_at`
	err := s.pool.QueryRow(ctx, query, p.Slug, p.Title, p.Tagline, p.Description, p.Tech,
		p.RepoURL, p.DemoURL, p.ImageURL, p.Featured, p.SortOrder, p.Published).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return &p, nil
}

func (s *Store) UpdateProject(ctx context.Context, id int64, p Project) (*Project, error) {
	const query = `UPDATE projects SET slug = $2, title = $3, tagline = $4, description = $5,
	                      tech = $6, repo_url = $7, demo_url = $8, image_url = $9,
	                      featured = $10, sort_order = $11, published = $12, updated_at = now()
	               WHERE id = $1
	               RETURNING created_at, updated_at`
	err := s.pool.QueryRow(ctx, query, id, p.Slug, p.Title, p.Tagline, p.Description, p.Tech,
		p.RepoURL, p.DemoURL, p.ImageURL, p.Featured, p.SortOrder, p.Published).
		Scan(&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	p.ID = id
	return &p, nil
}

func (s *Store) DeleteProject(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
