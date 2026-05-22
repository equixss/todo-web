package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User        UserDTOResponse `json:"user"`
	AccessToken string          `json:"access_token"`
	// Время истечения токена (Unix timestamp)
	ExpiresAt int64 `json:"expires_at"`
}

func LoginResponseFromDomain(result domain.LoginResult) LoginResponse {
	return LoginResponse{
		User:        UserDTOFromDomain(result.User),
		AccessToken: result.Tokens.AccessToken,
		ExpiresAt:   result.Tokens.ExpiresAt,
	}
}

// @Summary Вход в систему
// @Description Аутентификация пользователя по email или телефону. Возвращает JWT токен доступа и refresh токен.
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Запрос на вход"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string "Неверные учетные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/login [post]
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

	response := LoginResponseFromDomain(result)
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
