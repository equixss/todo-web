package core_http_middleware

import (
	"time"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("X-Request-ID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
			c.Set("X-Request-ID", requestID)
			c.Header("X-Request-ID", requestID)
		}
		l := log.With(
			zap.String("request_id", requestID),
			zap.String("url", c.Request.URL.String()),
		)

		ctx := core_logger.ToContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Trace() Middleware {
	return func(c *gin.Context) {
		log := core_logger.FromContext(c.Request.Context())
		before := time.Now()
		log.Debug(
			">>>request started",
			zap.String("method", c.Request.Method),
			zap.Time("time", before.UTC()),
		)
		c.Next()

		log.Debug(
			"<<<request ended",
			zap.Int("statusCode", c.Writer.Status()),
			zap.Duration("latency", time.Since(before)),
		)
	}
}

func Panic(presenter *core_http_response.HTTPResponsePresenter) Middleware {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				presenter.PanicResponse(c, p, "HTTP unexpected panic")
			}
		}()
		c.Next()
	}
}
