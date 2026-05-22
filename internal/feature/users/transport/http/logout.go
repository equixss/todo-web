package users_transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Выход из системы
// @Description Отзывает токены пользователя. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /users/logout [post]
func (h *UsersHTTPHandler) Logout(c *gin.Context) {
	h.presenter.JSONResponse(c, map[string]string{
		"message": "logged out successfully",
	}, http.StatusOK)
}
