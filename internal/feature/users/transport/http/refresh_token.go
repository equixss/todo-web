package users_transport_http

import (
	"net/http"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Tokens AuthTokensDTO `json:"tokens"`
}

func (h *UsersHTTPHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke RefreshToken handler")

	var request RefreshTokenRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	result, err := h.usersService.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	response := RefreshTokenResponse{
		Tokens: AuthTokensFromDomain(result.Tokens),
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
