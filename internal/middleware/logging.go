package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// ZapLogger creates a middleware that logs HTTP requests using zap
func ZapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture the status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			defer func() {
				// Record the request completion
				duration := time.Since(start)

				// Get the request ID if it's available
				requestID := middleware.GetReqID(r.Context())

				// Prepare fields to log
				fields := []zap.Field{
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("request_id", requestID),
					zap.String("method", r.Method),
					zap.String("uri", r.RequestURI),
					zap.String("protocol", r.Proto),
					zap.Int("status", ww.Status()),
					zap.Int("bytes_written", ww.BytesWritten()),
					zap.Duration("duration", duration),
					zap.String("user_agent", r.UserAgent()),
					zap.String("referer", r.Referer()),
				}

				// Log at appropriate level based on status code
				statusCode := ww.Status()
				if statusCode >= 500 {
					logger.Error("Server error", fields...)
				} else if statusCode >= 400 {
					logger.Warn("Client error", fields...)
				} else {
					logger.Info("Request completed", fields...)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
