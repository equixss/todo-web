package users_transport_http

import (
	"net/http"

	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Tokens AuthTokensDTO `json:"tokens"`
}

func (h *UsersHTTPHandler) RefreshToken(c *gin.Context) {
	var request RefreshTokenRequest

	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate HTTP request")
		return
	}

	result, err := h.usersService.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to refresh token")
		return
	}

	response := RefreshTokenResponse{
		Tokens: AuthTokensFromDomain(result.Tokens),
	}
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
