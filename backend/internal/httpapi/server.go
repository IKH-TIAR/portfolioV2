package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type UserStore interface {
	FindUserByEmail(ctx context.Context, email string) (*store.User, error)
}

var _ UserStore = (*store.Store)(nil)

type Pinger interface {
	Ping(ctx context.Context) error
}

var _ Pinger = (*store.Store)(nil)

type Server struct {
	logger  *slog.Logger
	users   UserStore
	tokens  *auth.TokenService
	pinger  Pinger
	limiter LoginLimiter
	origins map[string]bool
}

type Options struct {
	Logger  *slog.Logger
	Users   UserStore
	Tokens  *auth.TokenService
	Pinger  Pinger
	Limiter LoginLimiter
	Origins []string
}

func New(opts Options) *Server {
	origins := make(map[string]bool, len(opts.Origins))
	for _, o := range opts.Origins {
		origins[o] = true
	}
	return &Server{
		logger:  opts.Logger,
		users:   opts.Users,
		tokens:  opts.Tokens,
		pinger:  opts.Pinger,
		limiter: opts.Limiter,
		origins: origins,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.requireAuth(s.handleLogout))
	mux.HandleFunc("GET /api/auth/me", s.requireAuth(s.handleMe))

	handler := recoverMiddleware(s.logger, mux)
	handler = loggingMiddleware(s.logger, handler)
	handler = corsMiddleware(s.origins, handler)
	return handler
}
