package subscriptionservice

import (
	"net/http"
	"time"

	"efmob/logger"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// requestLogMiddleware пишет в лог входящий запрос и итог после ответа (статус, длительность).
func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Log().Info("http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("raw_query", r.URL.RawQuery),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.Header.Get("User-Agent")),
			zap.String("content_type", r.Header.Get("Content-Type")),
			zap.Int64("content_length", r.ContentLength),
		)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		logger.Log().Info("http response",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Int("bytes", ww.BytesWritten()),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
