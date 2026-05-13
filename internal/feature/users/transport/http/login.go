package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Login handler")

	var request LoginRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	credentials := domain.NewLoginCredentials(request.Identifier, request.Password)

	result, err := h.usersService.Login(ctx, credentials)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	response := LoginResponseFromDomain(result)

	responseHandler.JSONResponse(response, http.StatusOK)
}
