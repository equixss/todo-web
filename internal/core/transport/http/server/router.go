package core_http_server

import (
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
	APIVersion2 = APIVersion("v2")
)

type APIVersionRouter struct {
	*gin.RouterGroup
	apiVersion APIVersion
	Middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
	gin *gin.Engine,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	ginGroup := gin.Group("api/" + string(apiVersion))

	return &APIVersionRouter{
		RouterGroup: ginGroup,
		apiVersion:  apiVersion,
		Middleware:  middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		handlers := make([]gin.HandlerFunc, 0, len(route.Middleware)+1)

		for _, m := range route.Middleware {
			handlers = append(handlers, gin.HandlerFunc(m))
		}
		handlers = append(handlers, route.Handler)
		r.Handle(route.Method, route.Path, handlers...)
	}
}
