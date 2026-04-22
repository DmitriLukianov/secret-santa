package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	imgopt "secret-santa-backend/internal/image"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"

	"github.com/google/uuid"
)

const maxUploadSize = 5 << 20 // 5 MB

// FileStorage is satisfied by *storage.S3.
type FileStorage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
}

type UploadHandler struct {
	s3 FileStorage
}

func NewUploadHandler(s3 FileStorage) *UploadHandler {
	return &UploadHandler{s3: s3}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	contentType := http.DetectContentType(raw)
	if !strings.HasPrefix(contentType, "image/") {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	optimized, mimeType, optErr := imgopt.Optimize(raw)
	if optErr != nil {
		optimized = raw
		mimeType = contentType
	}

	key := uuid.New().String() + ".jpg"
	url, err := h.s3.Upload(r.Context(), key, optimized, mimeType)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
