package users_transport_http

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_utils "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type DeleteUserResponse UserDTOResponse

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	requestedUserID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get ID path param")
		return
	}

	authenticatedUserID, ok := core_http_middleware.GetUserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("user not authenticated"), "authentication required")
		return
	}

	if authenticatedUserID != requestedUserID {
		responseHandler.ErrorResponse(fmt.Errorf("access denied"), "access denied")
		return
	}

	if err := h.usersService.DeleteUser(ctx, requestedUserID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}
	responseHandler.ResponseNoContent()
}
