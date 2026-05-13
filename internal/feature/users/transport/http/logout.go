package users_transport_http

import (
	"net/http"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Logout handler")

	responseHandler.JSONResponse(map[string]string{
		"message": "logged out successfully",
	}, http.StatusOK)
}