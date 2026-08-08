package httpapi

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ikhtiar/portfolio-backend/internal/store"
	"github.com/jackc/pgx/v5/pgconn"
)

func contextWithTimeout(ctx context.Context, seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
}

func htmlEscape(s string) string {
	return html.EscapeString(s)
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func parseID(r *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return id, nil
}

func paginationFrom(r *http.Request) (limit, offset int) {
	limit = 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

func isValidSlug(s string) bool {
	return len(s) >= 1 && len(s) <= 100 && slugPattern.MatchString(s)
}

func isValidKey(s string) bool {
	return len(s) >= 1 && len(s) <= 64 && slugPattern.MatchString(s)
}

func isValidURL(s string) bool {
	return s == "" || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func isValidEmail(s string) bool {
	return strings.Contains(s, "@") && len(s) <= 254 && len(s) >= 3
}

func storeError(logger *slog.Logger, w http.ResponseWriter, err error, action string) {
	var pgErr *pgconn.PgError
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		writeError(w, http.StatusConflict, "a record with that slug already exists")
		return
	}
	logger.Error(action, "error", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}
