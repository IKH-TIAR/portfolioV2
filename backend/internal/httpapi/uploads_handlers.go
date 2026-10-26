package httpapi

import (
	"crypto/rand"
	"encoding/hex"
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
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	if len(files) > 20 {
		writeError(w, http.StatusBadRequest, "at most 20 files per request")
		return
	}

	now := time.Now()
	results := make([]map[string]any, 0, len(files))
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			writeError(w, http.StatusBadRequest, "could not read file")
			return
		}
		contentType := header.Header.Get("Content-Type")
		ext, ok := allowedImageTypes[contentType]
		if !ok {
			file.Close()
			writeError(w, http.StatusBadRequest, "only jpeg, png, gif, webp and avif images are allowed")
			return
		}

		data, err := io.ReadAll(io.LimitReader(file, maxUploadBytes))
		file.Close()
		if err != nil {
			writeError(w, http.StatusBadRequest, "could not read file")
			return
		}

		objectPath := fmt.Sprintf("uploads/%04d/%02d/%s.%s", now.Year(), now.Month(), randomHex(16), ext)
		url, err := s.uploader.Upload(r.Context(), objectPath, data, contentType)
		if err != nil {
			s.logger.Error("upload image", "error", err)
			writeError(w, http.StatusInternalServerError, "upload failed")
			return
		}
		results = append(results, map[string]any{
			"url":         url,
			"contentType": contentType,
			"size":        len(data),
			"path":        objectPath,
		})
	}

	if len(results) == 1 {
		writeJSON(w, http.StatusOK, results[0])
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"urls": results})
}

func randomHex(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
