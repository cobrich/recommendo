package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func NewLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Info("Request started", "method", r.Method, "path", r.URL.Path)

			// Вызываем следующий хендлер в цепочке
			next.ServeHTTP(w, r)

			// После того как хендлер отработал, логируем информацию о запросе
			logger.Info("Handled request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}
