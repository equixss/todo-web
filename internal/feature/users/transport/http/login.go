package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

func (h *UsersHTTPHandler) Login(c *gin.Context) {
	var request LoginRequest

	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate HTTP request")
		return
	}

	credentials := domain.NewLoginCredentials(request.Identifier, request.Password)

	result, err := h.usersService.Login(c.Request.Context(), credentials)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to login")
		return
	}

	response := LoginResponse{
		User:        UserDTOFromDomain(result.User),
		AccessToken: result.Tokens.AccessToken,
		ExpiresAt:   result.Tokens.ExpiresAt,
	}
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
