package core_http_server

import (
	"github.com/gin-gonic/gin"

	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    gin.HandlerFunc
	Middleware []core_http_middleware.Middleware
}
