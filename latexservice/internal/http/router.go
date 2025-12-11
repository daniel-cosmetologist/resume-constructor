package http

import (
	"log"
	stdhttp "net/http"

	"latex_service/internal/latex"
)

// Server инкапсулирует HTTP-маршрутизацию latex-service.
type Server struct {
	mux      *stdhttp.ServeMux
	renderer *latex.Renderer
	logger   *log.Logger
}

// NewServer создаёт новый HTTP-сервер latex-service.
func NewServer(renderer *latex.Renderer, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.Default()
	}

	mux := stdhttp.NewServeMux()
	s := &Server{
		mux:      mux,
		renderer: renderer,
		logger:   logger,
	}

	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/internal/v1/render", s.handleRender)

	return s
}

// ServeHTTP реализует http.Handler.
func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s.mux.ServeHTTP(w, r)
}
