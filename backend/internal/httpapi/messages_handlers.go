package httpapi

import (
	"fmt"
	"net/http"
	"strings"
)

type contactRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Body  string `json:"body"`
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	if !s.contactLimiter.Allow(clientIP(r)) {
		writeError(w, http.StatusTooManyRequests, "too many messages, try again later")
		return
	}

	var req contactRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Body = strings.TrimSpace(req.Body)
	if req.Name == "" || len(req.Name) > 120 {
		writeError(w, http.StatusBadRequest, "name is required (max 120 chars)")
		return
	}
	if !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	if req.Body == "" || len(req.Body) > 5000 {
		writeError(w, http.StatusBadRequest, "message is required (max 5000 chars)")
		return
	}

	message, err := s.store.CreateMessage(r.Context(), req.Name, req.Email, req.Body)
	if err != nil {
		storeError(s.logger, w, err, "create message")
		return
	}

	if s.mailer != nil && s.adminEmail != "" {
		subject := fmt.Sprintf("New portfolio message from %s", req.Name)
		html := fmt.Sprintf(
			"<p><strong>Name:</strong> %s</p><p><strong>Email:</strong> %s</p><p><strong>Message:</strong></p><p>%s</p>",
			htmlEscape(req.Name), htmlEscape(req.Email), htmlEscape(req.Body))
		go func() {
			ctx, cancel := contextWithTimeout(r.Context(), 15)
			defer cancel()
			if err := s.mailer.Send(ctx, s.adminEmail, subject, html); err != nil {
				s.logger.Error("send contact email", "error", err)
			}
		}()
	}

	writeJSON(w, http.StatusCreated, message)
}

func (s *Server) handleAdminListMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := s.store.ListMessages(r.Context())
	if err != nil {
		storeError(s.logger, w, err, "list messages")
		return
	}
	writeJSON(w, http.StatusOK, messages)
}

func (s *Server) handleAdminDeleteMessage(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeleteMessage(r.Context(), id); err != nil {
		storeError(s.logger, w, err, "delete message")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
