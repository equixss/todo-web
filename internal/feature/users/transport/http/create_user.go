package users_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	users_service "github.com/equixss/todo-web/internal/feature/users/service"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Phone    *string `json:"phone" validate:"omitempty,min=10,max=15"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password string  `json:"password" validate:"required,min=6"`
}

type CreateUserResponse UserDTOResponse

// @Summary Создание нового пользователя
// @Description Создает нового пользователя в системе. Требуется авторизация.
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Запрос на создание пользователя"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} map[string]string "Некорректные данные запроса"
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *UsersHTTPHandler) CreateUser(c *gin.Context) {
	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate HTTP request")
		return
	}

	passwordHash, err := users_service.HashPassword(request.Password)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to hash password")
		return
	}

	userDomain, err := h.usersService.CreateUser(c.Request.Context(), domainFromDTO(request, passwordHash))

	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to create user")
		return
	}

	response := CreateUserResponse(UserDTOFromDomain(userDomain))

	h.presenter.JSONResponse(c, response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest, passwordHash string) domain.User {
	return domain.NewUserUninitialized(dto.Name, dto.Phone, dto.Email, passwordHash)
}
