package http

import (
	"log"
	stdhttp "net/http"
	"strings"
	"time"
)

// Middleware описывает функцию-обёртку над http.Handler.
type Middleware func(stdhttp.Handler) stdhttp.Handler)

// applyMiddleware последовательно оборачивает обработчик в цепочку middleware.
func (s *Server) applyMiddleware(h stdhttp.Handler, mws ...Middleware) stdhttp.Handler {
	if len(mws) == 0 {
		return h
	}
	// Применяем в обратном порядке: последний в списке выполняется первым.
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// LoggingMiddleware логирует входящие запросы и статус-ответ.
func LoggingMiddleware(logger *log.Logger) Middleware {
	if logger == nil {
		logger = log.Default()
	}

	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			start := time.Now()

			// Обёртка для захвата статуса.
			ww := &responseWriterWrapper{ResponseWriter: w, status: stdhttp.StatusOK}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			logger.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.status, duration)
		})
	}
}

// RecoverMiddleware перехватывает паники и возвращает 500 вместо падения процесса.
func RecoverMiddleware(logger *log.Logger) Middleware {
	if logger == nil {
		logger = log.Default()
	}

	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Printf("panic recovered: %v", rec)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(stdhttp.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal_error","message":"Internal server error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// JSONOnlyMiddleware гарантирует, что запросы имеют Content-Type application/json.
func JSONOnlyMiddleware() Middleware {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			if r.Method == stdhttp.MethodPost || r.Method == stdhttp.MethodPut || r.Method == stdhttp.MethodPatch {
				ct := r.Header.Get("Content-Type")
				if ct == "" || !strings.HasPrefix(strings.ToLower(ct), "application/json") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(stdhttp.StatusUnsupportedMediaType)
					_, _ = w.Write([]byte(`{"error":"unsupported_media_type","message":"Content-Type must be application/json"}`))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriterWrapper используется для захвата статуса ответа.
type responseWriterWrapper struct {
	stdhttp.ResponseWriter
	status int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
