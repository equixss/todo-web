package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_server "github.com/equixss/todo-web/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
	jwtMW             *core_http_middleware.JWTMiddleware
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
	jwtMW *core_http_middleware.JWTMiddleware,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
		jwtMW:             jwtMW,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/statistics",
			Handler:    h.GetStatistics,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
	}
}
