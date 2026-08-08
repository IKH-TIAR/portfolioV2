package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
)

type contextKey string

const authContextKey contextKey = "auth"

var ErrUnauthorized = errors.New("unauthorized")

func claimsFrom(r *http.Request) *auth.Claims {
	claims, _ := r.Context().Value(authContextKey).(*auth.Claims)
	return claims
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || token == "" {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		claims, err := s.tokens.Parse(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), authContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}
