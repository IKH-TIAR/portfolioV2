package httpapi

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow(clientIP(r)) {
		writeError(w, http.StatusTooManyRequests, "too many attempts, try again later")
		return
	}

	var req loginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := s.store.FindUserByEmail(r.Context(), req.Email)
	if errors.Is(err, store.ErrNotFound) || (err == nil && !auth.VerifyPassword(user.PasswordHash, req.Password)) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		s.logger.Error("login: lookup user", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	token, err := s.tokens.Issue(user.ID, user.Email)
	if err != nil {
		s.logger.Error("login: issue token", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"expires_in": int(s.tokens.MaxAge().Seconds()),
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	claims := claimsFrom(r)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":    claims.UserID,
		"email": claims.Email,
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
