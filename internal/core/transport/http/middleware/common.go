package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Header.Get("X-Request-ID")
				if requestID == "" {
					requestID = uuid.NewString()
				}
				r.Header.Set("X-Request-ID", requestID)
				w.Header().Set("X-Request-ID", requestID)

				next.ServeHTTP(w, r)
			})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			crw := core_http_response.NewResponseWriter(rw)

			before := time.Now()
			log.Debug(
				">>>request started",
				zap.String("method", r.Method),
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(crw, r)

			log.Debug(
				"<<<request ended",
				zap.Int("statusCode", crw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			handler := core_http_response.NewHTTPResponseHandler(log, rw)
			defer func() {
				if p := recover(); p != nil {
					handler.PanicResponse(p, "HTTP unexpected panic")
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
