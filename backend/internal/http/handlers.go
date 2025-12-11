package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"time"

	"resume_backend/internal/resume"
)

// handleHealth — простой health-check эндпоинт.
func (s *Server) handleHealth(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodGet && r.Method != stdhttp.MethodHead {
		w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"status":  "ok",
		"service": "resume-backend",
		"time":    time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// handleGeneratePDF принимает JSON-данные резюме, вызывает доменный сервис
// и возвращает PDF-файл, сгенерированный latex-service.
func (s *Server) handleGeneratePDF(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "method_not_allowed",
			"message": "Only POST is allowed",
		})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req resume.Resume

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_json",
			"message": fmt.Sprintf("Failed to parse request: %v", err),
		})
		return
	}

	pdf, err := s.resumeService.GeneratePDF(ctx, req)
	if err != nil {
		var ve *resume.ValidationError
		if errors.As(err, &ve) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(stdhttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":   "validation_error",
				"message": "Invalid resume data",
				"details": ve.Errors,
			})
			return
		}

		s.logger.Printf("GeneratePDF error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "generation_failed",
			"message": "Failed to generate PDF",
		})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.WriteHeader(stdhttp.StatusOK)
	_, _ = w.Write(pdf)
}
