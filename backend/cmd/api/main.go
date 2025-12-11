package main

import (
	"log"
	"net/http"

	"resume_backend/internal/config"
	httptransport "resume_backend/internal/http"
	"resume_backend/internal/latexclient"
	"resume_backend/internal/resume"
)

func main() {
	cfg := config.Load()

	logger := log.Default()
	logger.Printf(
		"starting resume-backend on %s (latex-service: %s)",
		cfg.HTTPAddr,
		cfg.LaTeXServiceURL,
	)

	// HTTP-клиент к latex-service
	latexClient := latexclient.NewClient(cfg.LaTeXServiceURL, logger)

	// Доменный сервис резюме, который валидирует данные и зовёт latex-service
	resumeService := resume.NewService(latexClient, logger)

	// HTTP-слой (REST API)
	server := httptransport.NewServer(resumeService, logger)

	if err := http.ListenAndServe(cfg.HTTPAddr, server); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
