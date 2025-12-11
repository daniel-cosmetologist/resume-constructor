package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"time"
)

// Максимальный размер тела запроса с резюме (примерное значение).
const maxResumeBodySize = 256 * 1024 // 256 KB

// ResumeService описывает зависимость уровня бизнес-логики для генерации PDF.
type ResumeService interface {
	GeneratePDF(ctx context.Context, resume ResumeRequest) ([]byte, error)
}

// FieldError описывает ошибку с полями валидации.
type FieldError interface {
	error
	Fields() map[string]string
}

// ResumeRequest описывает JSON-полезную нагрузку запроса.
type ResumeRequest struct {
	FullName       string            `json:"fullName"`
	Position       string            `json:"position"`
	Summary        string            `json:"summary"`
	Contacts       Contacts          `json:"contacts"`
	Skills         []string          `json:"skills"`
	Experience     []ExperienceEntry `json:"experience"`
	Education      []EducationEntry  `json:"education"`
	CustomSections []CustomSection   `json:"customSections"`
	Photo          *Photo            `json:"photo,omitempty"`
}

// Contacts описывает контактные данные пользователя.
type Contacts struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Links    []Link `json:"links"`
}

// Link описывает внешнюю ссылку (GitHub, LinkedIn и т.п.).
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ExperienceEntry описывает запись об опыте работы.
type ExperienceEntry struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Bullets     []string `json:"bullets"`
}

// EducationEntry описывает запись об образовании.
type EducationEntry struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Details     string `json:"details"`
}

// CustomSection описывает кастомную секцию с буллетами.
type CustomSection struct {
	Title        string   `json:"title"`
	BulletSymbol string   `json:"bulletSymbol"`
	Items        []string `json:"items"`
}

// Photo описывает загруженное фото в base64.
type Photo struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// errorResponse описывает формат JSON-ошибки, возвращаемой API.
type errorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// handleHealth простой health-check эндпоинт.
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
		"service": "resume-backend",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// handleGenerateResumePDF обрабатывает POST /api/v1/resume/pdf.
func (s *Server) handleGenerateResumePDF(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		writeJSONError(w, stdhttp.StatusMethodNotAllowed, "method_not_allowed", "Only POST is allowed", nil)
		return
	}

	if s.resumeService == nil {
		writeJSONError(w, stdhttp.StatusInternalServerError, "internal_error", "Resume service is not configured", nil)
		return
	}

	// Ограничиваем размер тела запроса.
	r.Body = stdhttp.MaxBytesReader(w, r.Body, maxResumeBodySize)
	defer r.Body.Close()

	var payload ResumeRequest
	if err := decodeJSON(r.Body, &payload); err != nil {
		writeJSONError(w, stdhttp.StatusBadRequest, "invalid_json", err.Error(), nil)
		return
	}

	ctx := r.Context()
	pdf, err := s.resumeService.GeneratePDF(ctx, payload)
	if err != nil {
		// Пытаемся извлечь ошибки валидации, если сервис их предоставляет.
		var ferr FieldError
		if errors.As(err, &ferr) {
			writeJSONError(w, stdhttp.StatusBadRequest, "validation_error", "Invalid resume data", ferr.Fields())
			return
		}

		s.logger.Printf("failed to generate PDF: %v", err)
		writeJSONError(w, stdhttp.StatusInternalServerError, "pdf_generation_failed", "Failed to generate PDF", nil)
		return
	}

	// Отдаём PDF.
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="resume.pdf"`)
	w.WriteHeader(stdhttp.StatusOK)

	if _, err := w.Write(pdf); err != nil {
		s.logger.Printf("failed to write PDF response: %v", err)
	}
}

// decodeJSON декодирует JSON с защитой от избыточных полей/размеров.
func decodeJSON(body io.Reader, dst any) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Проверяем, что в потоке нет лишних данных после объекта.
	if dec.More() {
		return fmt.Errorf("unexpected data after JSON payload")
	}

	return nil
}

// writeJSONError отправляет JSON-ошибку в стандартном формате.
func writeJSONError(w stdhttp.ResponseWriter, status int, code, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{
		Error:   code,
		Message: message,
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
