package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(r.Handler, r.Middleware...)
}
