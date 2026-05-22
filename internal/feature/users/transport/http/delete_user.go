package users_transport_http

import (
	"fmt"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_utils "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

// @Summary Удаление пользователя
// @Description Удаляет пользователя по его ID. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 "Пользователь успешно удален"
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(c *gin.Context) {
	requestedUserID, err := core_http_utils.GetIntPathValue(c.Request, "id")
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get ID path param")
		return
	}

	authenticatedUserID, ok := core_http_middleware.GetUserIDFromContext(c)
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	if authenticatedUserID != requestedUserID {
		h.presenter.ErrorResponse(c, fmt.Errorf("access denied"), "access denied")
		return
	}

	if err := h.usersService.DeleteUser(c, requestedUserID); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to delete user")
		return
	}
	h.presenter.ResponseNoContent(c)
}
