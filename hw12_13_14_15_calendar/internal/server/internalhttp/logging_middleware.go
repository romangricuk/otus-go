package internalhttp

import (
	"net/http"
	"time"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK, 0}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggingMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			clientIP := r.RemoteAddr
			method := r.Method
			path := r.URL.Path
			query := r.URL.RawQuery
			if query != "" {
				path = path + "?" + query
			}
			httpVersion := r.Proto
			statusCode := rw.statusCode
			size := rw.size
			userAgent := r.UserAgent()
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			logger.Infof(
				"%s [%s] %s %s %s %d %d \"%s\" %s %s",
				clientIP,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				method,
				path,
				httpVersion,
				statusCode,
				size,
				userAgent,
				latency,
				requestID,
			)
		})
	}
}
