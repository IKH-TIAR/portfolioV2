package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/ikhtiar/portfolio-backend/internal/store"
)

type photoPayload struct {
	URL       string     `json:"url"`
	Caption   string     `json:"caption"`
	Location  string     `json:"location"`
	TakenAt   *time.Time `json:"takenAt,omitempty"`
	Published bool       `json:"published"`
}

func (p *photoPayload) validate() string {
	if !isValidURL(p.URL) {
		return "url must be an http(s) URL"
	}
	if len(p.Caption) > 300 {
		return "caption must be at most 300 chars"
	}
	if len(p.Location) > 200 {
		return "location must be at most 200 chars"
	}
	return ""
}

func toPhoto(p photoPayload) store.Photo {
	return store.Photo{
		URL:       p.URL,
		Caption:   p.Caption,
		Location:  p.Location,
		TakenAt:   p.TakenAt,
		Published: p.Published,
	}
}

func (s *Server) handleListPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := s.store.ListPhotos(r.Context(), true)
	if err != nil {
		storeError(s.logger, w, err, "list photos")
		return
	}
	writeJSON(w, http.StatusOK, photos)
}

func (s *Server) handleAdminListPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := s.store.ListPhotos(r.Context(), false)
	if err != nil {
		storeError(s.logger, w, err, "list photos")
		return
	}
	writeJSON(w, http.StatusOK, photos)
}

func (s *Server) handleAdminGetPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	photo, err := s.store.GetPhotoByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "photo not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "get photo")
		return
	}
	writeJSON(w, http.StatusOK, photo)
}

func (s *Server) handleAdminCreatePhoto(w http.ResponseWriter, r *http.Request) {
	var payload photoPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	photo, err := s.store.CreatePhoto(r.Context(), toPhoto(payload))
	if err != nil {
		storeError(s.logger, w, err, "create photo")
		return
	}
	writeJSON(w, http.StatusCreated, photo)
}

func (s *Server) handleAdminUpdatePhoto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var payload photoPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	photo, err := s.store.UpdatePhoto(r.Context(), id, toPhoto(payload))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "photo not found")
		return
	}
	if err != nil {
		storeError(s.logger, w, err, "update photo")
		return
	}
	writeJSON(w, http.StatusOK, photo)
}

func (s *Server) handleAdminDeletePhoto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeletePhoto(r.Context(), id); err != nil {
		storeError(s.logger, w, err, "delete photo")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
