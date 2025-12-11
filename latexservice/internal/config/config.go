package config

import "os"

// Config описывает конфигурацию latex-service.
type Config struct {
	// HTTPAddr — адрес, на котором слушает HTTP-сервер, например ":8081".
	HTTPAddr string
	// TemplatePath — путь к LaTeX-шаблону, например "templates/resume_template.tex".
	TemplatePath string
}

// Load загружает конфигурацию из переменных окружения с дефолтами.
func Load() Config {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8081"
	}

	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "templates/resume_template.tex"
	}

	return Config{
		HTTPAddr:     httpAddr,
		TemplatePath: templatePath,
	}
}
