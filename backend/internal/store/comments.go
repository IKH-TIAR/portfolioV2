package store

import (
	"context"
	"fmt"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Store) ListCommentsByPost(ctx context.Context, postID int64) ([]Comment, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, post_id, author, body, created_at
		 FROM comments WHERE post_id = $1 ORDER BY created_at ASC, id ASC`, postID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	out := make([]Comment, 0)
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) CreateComment(ctx context.Context, postID int64, author, body string) (*Comment, error) {
	const query = `INSERT INTO comments (post_id, author, body) VALUES ($1, $2, $3)
	               RETURNING id, post_id, author, body, created_at`
	var c Comment
	err := s.pool.QueryRow(ctx, query, postID, author, body).
		Scan(&c.ID, &c.PostID, &c.Author, &c.Body, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return &c, nil
}

func (s *Store) DeleteComment(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
