package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type Store interface {
	FindUserByEmail(ctx context.Context, email string) (*store.User, error)

	ListProjects(ctx context.Context, q store.ProjectQuery) ([]store.Project, error)
	GetProjectBySlug(ctx context.Context, slug string, publishedOnly bool) (*store.Project, error)
	GetProjectByID(ctx context.Context, id int64) (*store.Project, error)
	CreateProject(ctx context.Context, p store.Project) (*store.Project, error)
	UpdateProject(ctx context.Context, id int64, p store.Project) (*store.Project, error)
	DeleteProject(ctx context.Context, id int64) error

	ListPosts(ctx context.Context, q store.PostQuery) ([]store.Post, error)
	GetPostBySlug(ctx context.Context, slug string, publishedOnly bool) (*store.Post, error)
	GetPostByID(ctx context.Context, id int64) (*store.Post, error)
	CreatePost(ctx context.Context, p store.Post) (*store.Post, error)
	UpdatePost(ctx context.Context, id int64, p store.Post) (*store.Post, error)
	DeletePost(ctx context.Context, id int64) error

	ListPhotos(ctx context.Context, publishedOnly bool) ([]store.Photo, error)
	GetPhotoByID(ctx context.Context, id int64) (*store.Photo, error)
	CreatePhoto(ctx context.Context, p store.Photo) (*store.Photo, error)
	UpdatePhoto(ctx context.Context, id int64, p store.Photo) (*store.Photo, error)
	DeletePhoto(ctx context.Context, id int64) error

	ListSiteSettings(ctx context.Context) (map[string]json.RawMessage, error)
	SetSiteSetting(ctx context.Context, key string, value json.RawMessage) error

	ListCommentsByPost(ctx context.Context, postID int64) ([]store.Comment, error)
	CreateComment(ctx context.Context, postID int64, author, body string) (*store.Comment, error)
	DeleteComment(ctx context.Context, id int64) error

	CreateMessage(ctx context.Context, name, email, body string) (*store.Message, error)
	ListMessages(ctx context.Context) ([]store.Message, error)
	DeleteMessage(ctx context.Context, id int64) error
}

var _ Store = (*store.Store)(nil)

type Pinger interface {
	Ping(ctx context.Context) error
}

var _ Pinger = (*store.Store)(nil)

type Uploader interface {
	Upload(ctx context.Context, objectPath string, data []byte, contentType string) (string, error)
}

type Mailer interface {
	Send(ctx context.Context, to, subject, html string) error
}

type Server struct {
	logger         *slog.Logger
	store          Store
	tokens         *auth.TokenService
	pinger         Pinger
	limiter        LoginLimiter
	contactLimiter LoginLimiter
	uploader       Uploader
	mailer         Mailer
	adminEmail     string
	origins        map[string]bool
}

type Options struct {
	Logger         *slog.Logger
	Store          Store
	Tokens         *auth.TokenService
	Pinger         Pinger
	Limiter        LoginLimiter
	ContactLimiter LoginLimiter
	Uploader       Uploader
	Mailer         Mailer
	AdminEmail     string
	Origins        []string
}

func New(opts Options) *Server {
	origins := make(map[string]bool, len(opts.Origins))
	for _, o := range opts.Origins {
		origins[o] = true
	}
	return &Server{
		logger:         opts.Logger,
		store:          opts.Store,
		tokens:         opts.Tokens,
		pinger:         opts.Pinger,
		limiter:        opts.Limiter,
		contactLimiter: opts.ContactLimiter,
		uploader:       opts.Uploader,
		mailer:         opts.Mailer,
		adminEmail:     opts.AdminEmail,
		origins:        origins,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.requireAuth(s.handleLogout))
	mux.HandleFunc("GET /api/auth/me", s.requireAuth(s.handleMe))

	mux.HandleFunc("GET /api/site", s.handleGetSite)
	mux.HandleFunc("GET /api/projects", s.handleListProjects)
	mux.HandleFunc("GET /api/projects/{slug}", s.handleGetProject)
	mux.HandleFunc("GET /api/posts", s.handleListPosts)
	mux.HandleFunc("GET /api/posts/{slug}", s.handleGetPost)
	mux.HandleFunc("GET /api/posts/{slug}/comments", s.handleListComments)
	mux.HandleFunc("POST /api/posts/{slug}/comments", s.handleCreateComment)
	mux.HandleFunc("GET /api/photos", s.handleListPhotos)
	mux.HandleFunc("POST /api/contact", s.handleContact)

	mux.HandleFunc("GET /api/admin/projects", s.requireAuth(s.handleAdminListProjects))
	mux.HandleFunc("GET /api/admin/projects/{id}", s.requireAuth(s.handleAdminGetProject))
	mux.HandleFunc("POST /api/admin/projects", s.requireAuth(s.handleAdminCreateProject))
	mux.HandleFunc("PUT /api/admin/projects/{id}", s.requireAuth(s.handleAdminUpdateProject))
	mux.HandleFunc("DELETE /api/admin/projects/{id}", s.requireAuth(s.handleAdminDeleteProject))

	mux.HandleFunc("GET /api/admin/posts", s.requireAuth(s.handleAdminListPosts))
	mux.HandleFunc("GET /api/admin/posts/{id}", s.requireAuth(s.handleAdminGetPost))
	mux.HandleFunc("POST /api/admin/posts", s.requireAuth(s.handleAdminCreatePost))
	mux.HandleFunc("PUT /api/admin/posts/{id}", s.requireAuth(s.handleAdminUpdatePost))
	mux.HandleFunc("DELETE /api/admin/posts/{id}", s.requireAuth(s.handleAdminDeletePost))

	mux.HandleFunc("GET /api/admin/photos", s.requireAuth(s.handleAdminListPhotos))
	mux.HandleFunc("GET /api/admin/photos/{id}", s.requireAuth(s.handleAdminGetPhoto))
	mux.HandleFunc("POST /api/admin/photos", s.requireAuth(s.handleAdminCreatePhoto))
	mux.HandleFunc("PUT /api/admin/photos/{id}", s.requireAuth(s.handleAdminUpdatePhoto))
	mux.HandleFunc("DELETE /api/admin/photos/{id}", s.requireAuth(s.handleAdminDeletePhoto))

	mux.HandleFunc("GET /api/admin/site-settings", s.requireAuth(s.handleAdminGetSiteSettings))
	mux.HandleFunc("PUT /api/admin/site-settings/{key}", s.requireAuth(s.handleAdminSetSiteSetting))

	mux.HandleFunc("DELETE /api/admin/comments/{id}", s.requireAuth(s.handleAdminDeleteComment))

	mux.HandleFunc("GET /api/admin/messages", s.requireAuth(s.handleAdminListMessages))
	mux.HandleFunc("DELETE /api/admin/messages/{id}", s.requireAuth(s.handleAdminDeleteMessage))

	mux.HandleFunc("POST /api/admin/uploads", s.requireAuth(s.handleAdminUpload))

	handler := recoverMiddleware(s.logger, mux)
	handler = loggingMiddleware(s.logger, handler)
	handler = corsMiddleware(s.origins, handler)
	return handler
}
