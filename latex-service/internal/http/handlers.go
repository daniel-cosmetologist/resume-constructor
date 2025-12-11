package http

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"time"

	"latex_service/internal/model"
)

// errorResponse описывает формат JSON-ошибки latex-service.
type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const maxRenderBodySize = 256 * 1024 // 256 KB

// handleHealth обрабатывает /healthz.
func (s *Server) handleHealth(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodGet && r.Method != stdhttp.MethodHead {
		w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdhttp.StatusOK)

	resp := map[string]any{
		"status":  "ok",
		"time":    time.Now().UTC().Format(time.RFC3339),
		"service": "latex-service",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// handleRender обрабатывает POST /internal/v1/render.
func (s *Server) handleRender(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		writeJSONError(w, stdhttp.StatusMethodNotAllowed, "method_not_allowed", "Only POST is allowed")
		return
	}

	r.Body = stdhttp.MaxBytesReader(w, r.Body, maxRenderBodySize)
	defer r.Body.Close()

	var payload model.Resume
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&payload); err != nil {
		writeJSONError(w, stdhttp.StatusBadRequest, "invalid_json", fmt.Sprintf("Failed to decode JSON: %v", err))
		return
	}

	pdf, err := s.renderer.Render(r.Context(), payload)
	if err != nil {
		s.logger.Printf("failed to render PDF: %v", err)
		writeJSONError(w, stdhttp.StatusInternalServerError, "render_failed", "Failed to render PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(stdhttp.StatusOK)

	if _, err := w.Write(pdf); err != nil {
		s.logger.Printf("failed to write PDF response: %v", err)
	}
}

func writeJSONError(w stdhttp.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:   code,
		Message: msg,
	})
}
