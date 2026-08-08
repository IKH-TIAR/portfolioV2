package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type createCommentRequest struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (s *Server) handleListComments(w http.ResponseWriter, r *http.Request) {
	post, err := s.store.GetPostBySlug(r.Context(), r.PathValue("slug"), true)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get post for comments")
		return
	}
	comments, err := s.store.ListCommentsByPost(r.Context(), post.ID)
	if err != nil {
		storeError(s.logger, w, err, "list comments")
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (s *Server) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	post, err := s.store.GetPostBySlug(r.Context(), r.PathValue("slug"), true)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get post for comment")
		return
	}

	var req createCommentRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Author = strings.TrimSpace(req.Author)
	req.Body = strings.TrimSpace(req.Body)
	if req.Author == "" || len(req.Author) > 80 {
		writeError(w, http.StatusBadRequest, "author is required (max 80 chars)")
		return
	}
	if req.Body == "" || len(req.Body) > 2000 {
		writeError(w, http.StatusBadRequest, "body is required (max 2000 chars)")
		return
	}

	comment, err := s.store.CreateComment(r.Context(), post.ID, req.Author, req.Body)
	if err != nil {
		storeError(s.logger, w, err, "create comment")
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (s *Server) handleAdminDeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeleteComment(r.Context(), id); err != nil {
		storeError(s.logger, w, err, "delete comment")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
