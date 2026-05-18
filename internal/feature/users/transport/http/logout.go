package users_transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UsersHTTPHandler) Logout(c *gin.Context) {
	h.presenter.JSONResponse(c, map[string]string{
		"message": "logged out successfully",
	}, http.StatusOK)
}
