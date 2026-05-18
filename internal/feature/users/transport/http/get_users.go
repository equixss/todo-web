package users_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(c *gin.Context) {

	limit, offset, err := getLimitOffsetQueryParams(c.Request)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get limit/offset query params")
		return
	}
	usersDomain, err := h.usersService.GetUsers(c.Request.Context(), limit, offset)
	if err != nil {
		h.presenter.ErrorResponse(
			c,
			err,
			"failed to get users",
		)
		return
	}

	response := GetUsersResponse(UsersDTOFromDomains(usersDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf(`parameter "limit": %w`, err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf(`parameter "offset": %w`, err)
	}
	return limit, offset, nil
}
