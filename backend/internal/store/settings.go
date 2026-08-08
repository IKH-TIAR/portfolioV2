package store

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *Store) ListSiteSettings(ctx context.Context) (map[string]json.RawMessage, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, value FROM site_settings`)
	if err != nil {
		return nil, fmt.Errorf("list site settings: %w", err)
	}
	defer rows.Close()

	out := make(map[string]json.RawMessage)
	for rows.Next() {
		var key string
		var value json.RawMessage
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan site setting: %w", err)
		}
		out[key] = value
	}
	return out, rows.Err()
}

func (s *Store) SetSiteSetting(ctx context.Context, key string, value json.RawMessage) error {
	const query = `INSERT INTO site_settings (key, value) VALUES ($1, $2)
	               ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`
	if _, err := s.pool.Exec(ctx, query, key, value); err != nil {
		return fmt.Errorf("set site setting %q: %w", key, err)
	}
	return nil
}
