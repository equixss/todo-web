package users_transport_http

import (
	"net/http"

	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse AuthTokensDTO

// @Summary Обновление токена доступа
// @Description Обновляет JWT токен доступа используя refresh токен. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Запрос на обновление токена"
// @Success 200 {object} RefreshTokenResponse
// @Failure 401 {object} map[string]string "Требуется авторизация или недействительный refresh токен"
// @Router /users/refresh [post]
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

	response := RefreshTokenResponse(AuthTokensFromDomain(result.Tokens))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
