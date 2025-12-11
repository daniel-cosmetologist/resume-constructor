package config

import (
	"os"
)

// Config описывает конфигурацию backend-сервиса.
type Config struct {
	// HTTPAddr — адрес, на котором слушает HTTP-сервер, например ":8080".
	HTTPAddr string
	// LaTeXServiceURL — базовый URL latex-service, например "http://latex-service:8081".
	LaTeXServiceURL string
}

// Load загружает конфигурацию из переменных окружения с дефолтами.
func Load() Config {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	latexURL := os.Getenv("LATEX_SERVICE_URL")
	if latexURL == "" {
		latexURL = "http://latex-service:8081"
	}

	return Config{
		HTTPAddr:        httpAddr,
		LaTeXServiceURL: latexURL,
	}
}
