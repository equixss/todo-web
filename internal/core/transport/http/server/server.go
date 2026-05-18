package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPServer struct {
	Engine     *gin.Engine
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	if len(middleware) > 0 {
		engine.Use(middleware...)
	}

	return &HTTPServer{
		Engine:     engine,
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		if len(router.Middleware) > 0 {
			for _, m := range router.Middleware {
				router.RouterGroup.Use(m)
			}
		}
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		s.log.Warn("starting http server", zap.String("addr", s.config.Addr))
		err := s.Engine.Run(s.config.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutting down http server...")
		s.log.Warn("http server shutdown completed")
	}
	return nil
}

