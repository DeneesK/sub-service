package middlewares

import (
	"net/http"
	"time"
)

type Logger interface {
	Debugln(args ...interface{})
	Errorf(template string, args ...interface{})
	Error(args ...interface{})
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
		wroteHeader  bool
	}
)

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if !r.wroteHeader {
		r.ResponseWriter.WriteHeader(statusCode)
		r.responseData.status = statusCode
		r.wroteHeader = true
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func NewLoggingMiddleware(log Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			log.Debugln(
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		})
	}
}
