package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"resume_backend/internal/config"
	httptransport "resume_backend/internal/http"
)

// dummyResumeService — временная реализация ResumeService.
// На этапе интеграции с latex-service будет заменена реальной.
type dummyResumeService struct {
	latexServiceURL string
	logger          *log.Logger
}

func (s *dummyResumeService) GeneratePDF(ctx context.Context, resume httptransport.ResumeRequest) ([]byte, error) {
	// Здесь позднее будет реальный вызов latex-service.
	return nil, errors.New("GeneratePDF is not implemented yet")
}

func main() {
	cfg := config.Load()

	logger := log.Default()
	logger.Printf("starting resume-backend on %s", cfg.HTTPAddr)

	resumeSvc := &dummyResumeService{
		latexServiceURL: cfg.LaTeXServiceURL,
		logger:          logger,
	}

	server := httptransport.NewServer(resumeSvc, logger)

	if err := http.ListenAndServe(cfg.HTTPAddr, server); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
