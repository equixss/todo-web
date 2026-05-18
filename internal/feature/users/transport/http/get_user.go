package users_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(c *gin.Context) {
	requestedUserID, err := core_http_request.GetIntPathValue(c.Request, "id")
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get ID path param")
		return
	}

	authenticatedUserID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	if authenticatedUserID != requestedUserID {
		h.presenter.ErrorResponse(c, fmt.Errorf("access denied"), "access denied")
		return
	}

	userDomain, err := h.usersService.GetUser(c.Request.Context(), requestedUserID)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get user")
		return
	}
	response := GetUserResponse(UserDTOFromDomain(userDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
