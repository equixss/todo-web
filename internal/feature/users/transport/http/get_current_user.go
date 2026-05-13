package users_transport_http

import (
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) GetCurrentUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetCurrentUser handler")

	userID, ok := core_http_middleware.GetUserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			ErrUnauthorized,
			"failed to get current user",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get current user")
		return
	}

	response := UserDTOFromDomain(user)

	responseHandler.JSONResponse(response, http.StatusOK)
}

var ErrUnauthorized = core_errors.ErrUnauthorized
