package http

import (
	"context"
	"log"
	stdhttp "net/http"

	"resume_backend/internal/resume"
)

// ResumeService — интерфейс доменного сервиса, который знает,
// как из модели Resume сделать PDF.
type ResumeService interface {
	GeneratePDF(ctx context.Context, req resume.Resume) ([]byte, error)
}

// Server инкапсулирует HTTP-маршрутизацию backend API.
type Server struct {
	mux           *stdhttp.ServeMux
	resumeService ResumeService
	logger        *log.Logger
}

// NewServer создаёт новый экземпляр HTTP-сервера.
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

	s.registerRoutes()

	return s
}

// registerRoutes регистрирует маршруты и навешивает middleware.
func (s *Server) registerRoutes() {
	// health-check
	s.mux.Handle(
		"/healthz",
		s.applyMiddleware(
			stdhttp.HandlerFunc(s.handleHealth),
			LoggingMiddleware(s.logger),
			RecoverMiddleware(s.logger),
		),
	)

	// генерация PDF по данным резюме
	s.mux.Handle(
		"/api/v1/resume/pdf",
		s.applyMiddleware(
			stdhttp.HandlerFunc(s.handleGeneratePDF),
			LoggingMiddleware(s.logger),
			RecoverMiddleware(s.logger),
			JSONOnlyMiddleware(),
		),
	)
}

// ServeHTTP реализует интерфейс http.Handler.
func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s.mux.ServeHTTP(w, r)
}
