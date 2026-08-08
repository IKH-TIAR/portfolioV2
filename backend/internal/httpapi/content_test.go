package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

func TestAdminEndpointsRequireAuth(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})

	cases := []struct {
		method, path, body string
	}{
		{http.MethodGet, "/api/admin/projects", ""},
		{http.MethodPost, "/api/admin/projects", `{"slug":"x","title":"x"}`},
		{http.MethodPost, "/api/admin/posts", `{"slug":"x","title":"x"}`},
		{http.MethodPost, "/api/admin/photos", `{"url":"https://e.com/a.jpg"}`},
		{http.MethodGet, "/api/admin/site-settings", ""},
		{http.MethodPut, "/api/admin/site-settings/hero", `{"value":"hi"}`},
		{http.MethodGet, "/api/admin/messages", ""},
		{http.MethodPost, "/api/admin/uploads", ""},
		{http.MethodDelete, "/api/admin/comments/1", ""},
		{http.MethodDelete, "/api/admin/messages/1", ""},
	}
	for _, c := range cases {
		rec := doRequest(t, s, c.method, c.path, c.body, "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s: expected 401, got %d", c.method, c.path, rec.Code)
		}
	}
}

func TestCreateProjectValidation(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	token, err := s.tokens.Issue(1, "admin@test.dev")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name, body string
	}{
		{"bad slug", `{"slug":"Bad Slug!","title":"x"}`},
		{"missing title", `{"slug":"ok-slug"}`},
		{"bad url", `{"slug":"ok-slug","title":"x","repoUrl":"not-a-url"}`},
	}
	for _, c := range cases {
		rec := doRequest(t, s, http.MethodPost, "/api/admin/projects", c.body, token)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %d: %s", c.name, rec.Code, rec.Body.String())
		}
	}
}

func TestCreateCommentSuccess(t *testing.T) {
	st := &fakeStore{posts: map[string]*store.Post{
		"hello": {ID: 7, Slug: "hello", Title: "Hello", Published: true},
	}}
	s := newTestServer(t, st, &fakePinger{})

	rec := doRequest(t, s, http.MethodPost, "/api/posts/hello/comments",
		`{"author":"Dipto","body":"Nice post!"}`, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var body store.Comment
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Author != "Dipto" || body.PostID != 7 {
		t.Fatalf("unexpected comment: %+v", body)
	}
}

func TestCreateCommentValidation(t *testing.T) {
	st := &fakeStore{posts: map[string]*store.Post{
		"hello": {ID: 7, Slug: "hello", Title: "Hello", Published: true},
	}}
	s := newTestServer(t, st, &fakePinger{})

	rec := doRequest(t, s, http.MethodPost, "/api/posts/hello/comments",
		`{"author":"","body":"hi"}`, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty author, got %d", rec.Code)
	}

	rec = doRequest(t, s, http.MethodPost, "/api/posts/hello/comments",
		`{"author":"Dipto","body":""}`, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty body, got %d", rec.Code)
	}
}

func TestCommentOnMissingPost(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	rec := doRequest(t, s, http.MethodPost, "/api/posts/nope/comments",
		`{"author":"Dipto","body":"hi"}`, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestContactSuccess(t *testing.T) {
	st := &fakeStore{}
	s := newTestServer(t, st, &fakePinger{})

	rec := doRequest(t, s, http.MethodPost, "/api/contact",
		`{"name":"Alice","email":"alice@example.com","body":"Hello there"}`, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var body store.Message
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Name != "Alice" || body.Email != "alice@example.com" {
		t.Fatalf("unexpected message: %+v", body)
	}
}

func TestContactValidation(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})

	cases := []struct {
		name, body string
	}{
		{"empty name", `{"name":"","email":"a@b.com","body":"hi"}`},
		{"bad email", `{"name":"A","email":"not-an-email","body":"hi"}`},
		{"empty body", `{"name":"A","email":"a@b.com","body":""}`},
	}
	for _, c := range cases {
		rec := doRequest(t, s, http.MethodPost, "/api/contact", c.body, "")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %d", c.name, rec.Code)
		}
	}
}

func TestContactRateLimit(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	for i := 0; i < 5; i++ {
		rec := doRequest(t, s, http.MethodPost, "/api/contact",
			`{"name":"A","email":"a@b.com","body":"hi"}`, "")
		if rec.Code != http.StatusCreated {
			t.Fatalf("attempt %d: expected 201, got %d", i+1, rec.Code)
		}
	}
	rec := doRequest(t, s, http.MethodPost, "/api/contact",
		`{"name":"A","email":"a@b.com","body":"hi"}`, "")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after limit, got %d", rec.Code)
	}
}

func TestSetSiteSettingValidation(t *testing.T) {
	s := newTestServer(t, &fakeStore{}, &fakePinger{})
	token, err := s.tokens.Issue(1, "admin@test.dev")
	if err != nil {
		t.Fatal(err)
	}

	rec := doRequest(t, s, http.MethodPut, "/api/admin/site-settings/Bad%20Key!", `{"value":"x"}`, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid key, got %d", rec.Code)
	}
}

func TestListComments(t *testing.T) {
	st := &fakeStore{posts: map[string]*store.Post{
		"hello": {ID: 7, Slug: "hello", Title: "Hello", Published: true},
	}}
	s := newTestServer(t, st, &fakePinger{})
	rec := doRequest(t, s, http.MethodGet, "/api/posts/hello/comments", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
