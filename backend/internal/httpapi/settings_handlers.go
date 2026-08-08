package httpapi

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetSite(w http.ResponseWriter, r *http.Request) {
	settings, err := s.store.ListSiteSettings(r.Context())
	if err != nil {
		storeError(s.logger, w, err, "list site settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) handleAdminGetSiteSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.store.ListSiteSettings(r.Context())
	if err != nil {
		storeError(s.logger, w, err, "list site settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

type setSettingRequest struct {
	Value json.RawMessage `json:"value"`
}

func (s *Server) handleAdminSetSiteSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if !isValidKey(key) {
		writeError(w, http.StatusBadRequest, "invalid setting key")
		return
	}
	var req setSettingRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Value) == 0 {
		writeError(w, http.StatusBadRequest, "value is required")
		return
	}
	if !json.Valid(req.Value) {
		writeError(w, http.StatusBadRequest, "value must be valid JSON")
		return
	}
	if err := s.store.SetSiteSetting(r.Context(), key, req.Value); err != nil {
		storeError(s.logger, w, err, "set site setting")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"key": key, "value": req.Value})
}
