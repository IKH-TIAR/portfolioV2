package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type projectPayload struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Tagline     string   `json:"tagline"`
	Description string   `json:"description"`
	Tech        []string `json:"tech"`
	RepoURL     string   `json:"repoUrl"`
	DemoURL     string   `json:"demoUrl"`
	ImageURL    string   `json:"imageUrl"`
	Featured    bool     `json:"featured"`
	SortOrder   int      `json:"sortOrder"`
	Published   bool     `json:"published"`
}

func (p *projectPayload) validate() string {
	if !isValidSlug(p.Slug) {
		return "slug must be lowercase letters, numbers and hyphens (max 100 chars)"
	}
	if p.Title == "" || len(p.Title) > 200 {
		return "title is required (max 200 chars)"
	}
	if len(p.Tagline) > 300 {
		return "tagline must be at most 300 chars"
	}
	if len(p.Description) > 100000 {
		return "description must be at most 100000 chars"
	}
	if len(p.Tech) > 50 {
		return "tech list must be at most 50 entries"
	}
	for i := range p.Tech {
		p.Tech[i] = strings.TrimSpace(p.Tech[i])
		if len(p.Tech[i]) > 50 {
			return "each tech entry must be at most 50 chars"
		}
	}
	if !isValidURL(p.RepoURL) || !isValidURL(p.DemoURL) || !isValidURL(p.ImageURL) {
		return "urls must be http(s) URLs"
	}
	return ""
}

func toProject(p projectPayload) store.Project {
	tech := p.Tech
	if tech == nil {
		tech = []string{}
	}
	return store.Project{
		Slug:        p.Slug,
		Title:       p.Title,
		Tagline:     p.Tagline,
		Description: p.Description,
		Tech:        tech,
		RepoURL:     p.RepoURL,
		DemoURL:     p.DemoURL,
		ImageURL:    p.ImageURL,
		Featured:    p.Featured,
		SortOrder:   p.SortOrder,
		Published:   p.Published,
	}
}

func (s *Server) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.ListProjects(r.Context(), store.ProjectQuery{PublishedOnly: true})
	if err != nil {
		storeError(s.logger, w, err, "list projects")
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (s *Server) handleGetProject(w http.ResponseWriter, r *http.Request) {
	project, err := s.store.GetProjectBySlug(r.Context(), r.PathValue("slug"), true)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (s *Server) handleAdminListProjects(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginationFrom(r)
	projects, err := s.store.ListProjects(r.Context(), store.ProjectQuery{Limit: limit, Offset: offset})
	if err != nil {
		storeError(s.logger, w, err, "list projects")
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (s *Server) handleAdminGetProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, err := s.store.GetProjectByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (s *Server) handleAdminCreateProject(w http.ResponseWriter, r *http.Request) {
	var payload projectPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	project, err := s.store.CreateProject(r.Context(), toProject(payload))
	if err != nil {
		storeError(s.logger, w, err, "create project")
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (s *Server) handleAdminUpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var payload projectPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	project, err := s.store.UpdateProject(r.Context(), id, toProject(payload))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "update project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (s *Server) handleAdminDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeleteProject(r.Context(), id); err != nil {
		storeError(s.logger, w, err, "delete project")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
