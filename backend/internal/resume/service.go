package resume

import (
	"context"
	"fmt"
	"log"
)

// PDFRenderer описывает зависимость сервиса от внешнего LaTeX-сервиса.
type PDFRenderer interface {
	RenderResume(ctx context.Context, r Resume) ([]byte, error)
}

// Service реализует бизнес-логику генерации PDF.
type Service struct {
	renderer PDFRenderer
	logger   *log.Logger
}

// NewService создаёт новый сервис резюме.
func NewService(renderer PDFRenderer, logger *log.Logger) *Service {
	if logger == nil {
		logger = log.Default()
	}
	return &Service{
		renderer: renderer,
		logger:   logger,
	}
}

// GeneratePDF валидирует данные резюме и делегирует генерацию LaTeX-сервису.
func (s *Service) GeneratePDF(ctx context.Context, r Resume) ([]byte, error) {
	if err := ValidateResume(r); err != nil {
		// ошибки валидации пробрасываем наверх как есть
		return nil, err
	}

	pdf, err := s.renderer.RenderResume(ctx, r)
	if err != nil {
		s.logger.Printf("RenderResume error: %v", err)
		return nil, fmt.Errorf("latex render failed: %w", err)
	}

	return pdf, nil
}
