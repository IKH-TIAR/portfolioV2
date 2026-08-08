package migrate

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ikhtiar/portfolio-backend/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Version int64
	Name    string
	SQL     string
}

func List() ([]Migration, error) {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	var out []Migration
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		version, err := parseVersion(e.Name())
		if err != nil {
			return nil, err
		}
		sql, err := migrations.FS.ReadFile(e.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", e.Name(), err)
		}
		out = append(out, Migration{Version: version, Name: e.Name(), SQL: string(sql)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out, nil
}

func parseVersion(name string) (int64, error) {
	prefix := strings.SplitN(name, "_", 2)[0]
	v, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("migration %q: invalid version prefix %q", name, prefix)
	}
	return v, nil
}

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	const ensure = `CREATE TABLE IF NOT EXISTS schema_migrations (
		version    BIGINT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`
	if _, err := pool.Exec(ctx, ensure); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	migrations, err := List()
	if err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, pool)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			fmt.Printf("skip   %s (already applied)\n", m.Name)
			continue
		}
		if err := apply(ctx, pool, m); err != nil {
			return err
		}
		applied[m.Version] = true
		fmt.Printf("apply  %s\n", m.Name)
	}
	return nil
}

func appliedVersions(ctx context.Context, pool *pgxpool.Pool) (map[int64]bool, error) {
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("query schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan schema_migrations: %w", err)
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

func apply(ctx context.Context, pool *pgxpool.Pool, m Migration) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", m.Name, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, m.SQL); err != nil {
		return fmt.Errorf("exec migration %s: %w", m.Name, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.Version); err != nil {
		return fmt.Errorf("record migration %s: %w", m.Name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", m.Name, err)
	}
	return nil
}

func PoolFromURL(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
