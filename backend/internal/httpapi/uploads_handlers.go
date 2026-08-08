package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"time"
)

const maxUploadBytes = 10 << 20

var allowedImageTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/gif":  "gif",
	"image/webp": "webp",
	"image/avif": "avif",
}

func (s *Server) handleAdminUpload(w http.ResponseWriter, r *http.Request) {
	if s.uploader == nil {
		writeError(w, http.StatusServiceUnavailable, "uploads are not configured")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		writeError(w, http.StatusBadRequest, "only jpeg, png, gif, webp and avif images are allowed")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxUploadBytes))
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read file")
		return
	}

	now := time.Now()
	objectPath := fmt.Sprintf("uploads/%04d/%02d/%s.%s", now.Year(), now.Month(), randomHex(16), ext)

	url, err := s.uploader.Upload(r.Context(), objectPath, data, contentType)
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			writeError(w, http.StatusBadRequest, "file too large")
			return
		}
		s.logger.Error("upload image", "error", err)
		writeError(w, http.StatusInternalServerError, "upload failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"url":         url,
		"contentType": contentType,
		"size":        len(data),
		"path":        objectPath,
	})
}

func randomHex(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
