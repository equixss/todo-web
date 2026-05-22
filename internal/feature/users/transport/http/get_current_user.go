package users_transport_http

import (
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type GetCurrentUserResponse UserDTOResponse

// @Summary Получение текущего пользователя
// @Description Возвращает данные текущего аутентифицированного пользователя. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} GetCurrentUserResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /users/me [get]
func (h *UsersHTTPHandler) GetCurrentUser(c *gin.Context) {
	log := core_logger.FromContext(c)
	log.Debug("invoke GetCurrentUser handler")

	userID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(
			c,
			core_errors.ErrUnauthorized,
			"failed to get current user",
		)
		return
	}

	user, err := h.usersService.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get current user")
		return
	}

	response := GetCurrentUserResponse(UserDTOFromDomain(user))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
