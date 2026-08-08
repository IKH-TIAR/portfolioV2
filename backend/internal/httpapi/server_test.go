package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type fakeStore struct {
	*store.Store
	user  *store.User
	err   error
	posts map[string]*store.Post
}

func (f *fakeStore) FindUserByEmail(ctx context.Context, email string) (*store.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.user == nil || f.user.Email != email {
		return nil, store.ErrNotFound
	}
	return f.user, nil
}

func (f *fakeStore) GetPostBySlug(ctx context.Context, slug string, publishedOnly bool) (*store.Post, error) {
	if p, ok := f.posts[slug]; ok && (!publishedOnly || p.Published) {
		return p, nil
	}
	return nil, store.ErrNotFound
}

func (f *fakeStore) CreateComment(ctx context.Context, postID int64, author, body string) (*store.Comment, error) {
	return &store.Comment{ID: 1, PostID: postID, Author: author, Body: body, CreatedAt: time.Now()}, nil
}

func (f *fakeStore) ListCommentsByPost(ctx context.Context, postID int64) ([]store.Comment, error) {
	return nil, nil
}

func (f *fakeStore) CreateMessage(ctx context.Context, name, email, body string) (*store.Message, error) {
	return &store.Message{ID: 1, Name: name, Email: email, Body: body, CreatedAt: time.Now()}, nil
}

type fakePinger struct {
	err error
}

func (f *fakePinger) Ping(ctx context.Context) error {
	return f.err
}

type fakeMailer struct {
	sent int
}

func (f *fakeMailer) Send(ctx context.Context, to, subject, html string) error {
	f.sent++
	return nil
}

func newTestServer(t *testing.T, st *fakeStore, pinger *fakePinger) *Server {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	hash, err := auth.HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if st.user == nil {
		st.user = &store.User{ID: 1, Email: "admin@test.dev", PasswordHash: hash}
	}
	tokens := auth.NewTokenService("test-secret-that-is-long-enough", time.Hour)
	return New(Options{
		Logger:         logger,
		Store:          st,
		Tokens:         tokens,
		Pinger:         pinger,
		Limiter:        NewTokenBucketLimiter(),
		ContactLimiter: NewTokenBucketLimiterWith(5, 1.0/6.0),
		Mailer:         &fakeMailer{},
		AdminEmail:     "admin@test.dev",
		Origins:        []string{"http://localhost:3000"},
	})
}

func doRequest(t *testing.T, s *Server, method, path, body string, token string) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body == "" {
		buf = &bytes.Buffer{}
	} else {
		buf = bytes.NewBufferString(body)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	return rec
}

func TestHealthOK(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodGet, "/api/health", "", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" || body["db"] != "up" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestHealthDegraded(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{err: errors.New("db down")})
	rec := doRequest(t, s, http.MethodGet, "/api/health", "", "")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestLoginSuccess(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodPost, "/api/auth/login",
		`{"email":"admin@test.dev","password":"correct-password"}`, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Token string `json:"token"`
		User  struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Token == "" {
		t.Fatal("expected token in response")
	}
	if body.User.Email != "admin@test.dev" {
		t.Fatalf("unexpected user: %+v", body.User)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodPost, "/api/auth/login",
		`{"email":"admin@test.dev","password":"wrong"}`, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLoginUnknownUser(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodPost, "/api/auth/login",
		`{"email":"nobody@test.dev","password":"whatever123"}`, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestLoginBadRequest(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodPost, "/api/auth/login", `{"email":""}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMeRequiresAuth(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodGet, "/api/auth/me", "", "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMeWithToken(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	token, err := s.tokens.Issue(1, "admin@test.dev")
	if err != nil {
		t.Fatal(err)
	}
	rec := doRequest(t, s, http.MethodGet, "/api/auth/me", "", token)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["email"] != "admin@test.dev" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestMeWithInvalidToken(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodGet, "/api/auth/me", "", "not-a-real-token")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestLoginRateLimit(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	for i := 0; i < 5; i++ {
		rec := doRequest(t, s, http.MethodPost, "/api/auth/login",
			`{"email":"admin@test.dev","password":"wrong"}`, "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("attempt %d: expected 401, got %d", i+1, rec.Code)
		}
	}
	rec := doRequest(t, s, http.MethodPost, "/api/auth/login",
		`{"email":"admin@test.dev","password":"wrong"}`, "")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after limit, got %d", rec.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected CORS origin header, got %q", got)
	}
}

func TestCORSBlocksUnknownOrigin(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no CORS header for unknown origin, got %q", got)
	}
}
