package http

import (
	"log"
	stdhttp "net/http"
)

// Server инкапсулирует HTTP-маршрутизацию и зависимости уровня HTTP.
type Server struct {
	mux           *stdhttp.ServeMux
	resumeService ResumeService
	logger        *log.Logger
}

// NewServer создаёт новый экземпляр HTTP-сервера с настроенным роутером.
func NewServer(resumeService ResumeService, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.Default()
	}

	mux := stdhttp.NewServeMux()
	s := &Server{
		mux:           mux,
		resumeService: resumeService,
		logger:        logger,
	}

	// Health-check без сложных middleware.
	mux.HandleFunc("/healthz", s.handleHealth)

	// Основной API-эндпоинт с middleware-цепочкой.
	pdfHandler := stdhttp.HandlerFunc(s.handleGenerateResumePDF)
	pdfHandler = s.applyMiddleware(
		pdfHandler,
		LoggingMiddleware(logger),
		RecoverMiddleware(logger),
		JSONOnlyMiddleware(),
	)

	mux.Handle("/api/v1/resume/pdf", pdfHandler)

	return s
}

// ServeHTTP реализует интерфейс http.Handler и делегирует обработку ServeMux.
func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s.mux.ServeHTTP(w, r)
}
