package httpapi

import (
	"errors"
	"net/http"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type postPayload struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	PostType  string `json:"postType"`
	Published bool   `json:"published"`
}

func (p *postPayload) validate() string {
	if !isValidSlug(p.Slug) {
		return "slug must be lowercase letters, numbers and hyphens (max 100 chars)"
	}
	if p.Title == "" || len(p.Title) > 200 {
		return "title is required (max 200 chars)"
	}
	if len(p.Body) > 500000 {
		return "body must be at most 500000 chars"
	}
	if p.PostType == "" {
		p.PostType = "journal"
	}
	if p.PostType != "journal" && p.PostType != "note" {
		return "postType must be 'journal' or 'note'"
	}
	return ""
}

func toPost(p postPayload) store.Post {
	return store.Post{
		Slug:      p.Slug,
		Title:     p.Title,
		Body:      p.Body,
		PostType:  p.PostType,
		Published: p.Published,
	}
}

func (s *Server) handleListPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginationFrom(r)
	posts, err := s.store.ListPosts(r.Context(), store.PostQuery{
		PublishedOnly: true, Limit: limit, Offset: offset,
	})
	if err != nil {
		storeError(s.logger, w, err, "list posts")
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	post, err := s.store.GetPostBySlug(r.Context(), r.PathValue("slug"), true)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get post")
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (s *Server) handleAdminListPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginationFrom(r)
	posts, err := s.store.ListPosts(r.Context(), store.PostQuery{Limit: limit, Offset: offset})
	if err != nil {
		storeError(s.logger, w, err, "list posts")
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (s *Server) handleAdminGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := s.store.GetPostByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get post")
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (s *Server) handleAdminCreatePost(w http.ResponseWriter, r *http.Request) {
	var payload postPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	post, err := s.store.CreatePost(r.Context(), toPost(payload))
	if err != nil {
		storeError(s.logger, w, err, "create post")
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (s *Server) handleAdminUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var payload postPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	post, err := s.store.UpdatePost(r.Context(), id, toPost(payload))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "update post")
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (s *Server) handleAdminDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeletePost(r.Context(), id); err != nil {
		storeError(s.logger, w, err, "delete post")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
