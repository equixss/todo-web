package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_types "github.com/equixss/todo-web/internal/core/transport/http/types"
	"github.com/gin-gonic/gin"
)

type PatchUserRequest struct {
	Name  core_http_types.Nullable[string] `json:"name"`
	Phone core_http_types.Nullable[string] `json:"phone"`
	Email core_http_types.Nullable[string] `json:"email"`
}

type PatchUserSwaggerRequest struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
	Email *string `json:"email"`
}

type PatchUserResponse UserDTOResponse

// @Summary Обновление пользователя
// @Description Частичное обновление данных пользователя. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body PatchUserSwaggerRequest true "Запрос на обновление"
// @Success 200 {object} PatchUserResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(c *gin.Context) {
	requestedUserID, err := core_http_request.GetIntPathValue(c.Request, "id")
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get ID path param")
		return
	}

	authenticatedUserID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	if authenticatedUserID != requestedUserID {
		h.presenter.ErrorResponse(c, fmt.Errorf("access denied"), "access denied")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := h.usersService.PatchUser(c.Request.Context(), requestedUserID, userPatch)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to patch user")
		return
	}

	response := PatchUserResponse(UserDTOFromDomain(userDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}

func (r *PatchUserRequest) Validate() error {
	if r.Name.Set {
		if r.Name.Value == nil {
			return fmt.Errorf("'Name' can't be null")
		}
		nameLen := len([]rune(*r.Name.Value))
		if nameLen < 3 || nameLen > 100 {
			return fmt.Errorf("'Name' must be between 3 and 100 symbols")
		}
	}
	if r.Phone.Set && r.Phone.Value != nil {
		phoneLen := len([]rune(*r.Phone.Value))
		if phoneLen < 10 || !domain.RussianPhoneRE.MatchString(*r.Phone.Value) {
			return fmt.Errorf("invalid phone number: %s: %w", *r.Phone.Value, core_errors.ErrInvalidArgument)
		}
	}
	return nil
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Name.ToDomain(),
		request.Phone.ToDomain(),
		request.Email.ToDomain(),
	)
}
