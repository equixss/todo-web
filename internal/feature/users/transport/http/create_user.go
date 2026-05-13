package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
	users_service "github.com/equixss/todo-web/internal/feature/users/service"
)

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Phone    *string `json:"phone" validate:"omitempty,min=10,max=15"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password string  `json:"password" validate:"required,min=6"`
}

type CreateUserResponse UserDTOResponse

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

	passwordHash, err := users_service.HashPassword(request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to hash password")
		return
	}

	userDomain, err := h.usersService.CreateUser(ctx, domainFromDTO(request, passwordHash))

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(UserDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest, passwordHash string) domain.User {
	return domain.NewUserUninitialized(dto.Name, dto.Phone, dto.Email, passwordHash)
}
