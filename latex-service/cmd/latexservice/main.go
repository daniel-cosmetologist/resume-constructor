package main

import (
	"log"
	"net/http"

	"latex_service/internal/config"
	httphandler "latex_service/internal/http"
	"latex_service/internal/latex"
)

func main() {
	cfg := config.Load()

	logger := log.Default()
	logger.Printf("starting latex-service on %s", cfg.HTTPAddr)

	renderer := latex.NewRenderer(cfg.TemplatePath, logger)

	server := httphandler.NewServer(renderer, logger)

	if err := http.ListenAndServe(cfg.HTTPAddr, server); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
