package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/qwinkki/http-rest-taskmgr/internal/core/logger"
	"go.uber.org/zap"

	core_http_response "github.com/qwinkki/http-rest-taskmgr/internal/core/transport/http/response"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RequestID := r.Header.Get(requestIDHeader)
			if RequestID == "" {
				RequestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, RequestID)

			w.Header().Set(requestIDHeader, RequestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RequestID := r.Header.Get(requestIDHeader)

			log.With(
				zap.String("request_id", RequestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(*log, w)

			defer func() {
				if err := recover(); err != nil {
					responseHandler.Handle(err, "unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now().UTC()
			log.Debug(">>> incomming http request", zap.Time("time", before))

			next.ServeHTTP(rw, r)

			log.Debug(
				">>> done http requeest",
				zap.Int("status_code", rw.GetStatusCodeOrPanic()),
				zap.Duration("Latency", time.Since(before)))
		})
	}
}
