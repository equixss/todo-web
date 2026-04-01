package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
	core_http_types "github.com/equixss/todo-web/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Name  core_http_types.Nullable[string] `json:"name"`
	Phone core_http_types.Nullable[string] `json:"phone"`
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
	if r.Phone.Set {
		if r.Phone.Value == nil {
			phoneLen := len([]rune(*r.Phone.Value))
			if phoneLen < 10 || !domain.RussianPhoneRE.MatchString(*r.Phone.Value) {
				return fmt.Errorf("invalid phone number: %s: %w", *r.Phone.Value, core_errors.ErrInvalidArgument)
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke PatchUser handler")

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get ID path param")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(UserDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Name.ToDomain(),
		request.Phone.ToDomain(),
	)
}
