package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name  string  `json:"name"  validate:"required,min=3,max=100"`
	Phone *string `json:"phone" validate:"omitempty,min=10,max=15"`
}

type CreateUserResponse struct {
	ID      int     `json:"id"`
	Version int     `json:"version"`
	Name    string  `json:"name"`
	Phone   *string `json:"phone"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain, err := h.usersService.CreateUser(ctx, domainFromDTO(request))

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := dtoFromDomain(userDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.Name, dto.Phone)
}

func dtoFromDomain(user domain.User) CreateUserResponse {
	return CreateUserResponse{
		Name:    user.Name,
		Phone:   user.Phone,
		ID:      user.ID,
		Version: user.Version,
	}
}
