package users_transport_http

import (
	"context"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_response_presenter "github.com/equixss/todo-web/internal/core/transport/http/response"
	core_http_server "github.com/equixss/todo-web/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	presenter    *core_http_response_presenter.HTTPResponsePresenter
	usersService UsersService
	jwtMW        *core_http_middleware.JWTMiddleware
}
type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id int,
	) error
	Login(
		ctx context.Context,
		credentials domain.LoginCredentials,
	) (domain.LoginResult, error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (domain.RefreshTokenResult, error)
}

func NewUsersHttpHandler(
	usersService UsersService,
	jwtMW *core_http_middleware.JWTMiddleware,
	presenter *core_http_response_presenter.HTTPResponsePresenter,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{usersService: usersService, jwtMW: jwtMW, presenter: presenter}
}

func (h *UsersHTTPHandler) PublicRoutes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/refresh",
			Handler: h.RefreshToken,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}

func (h *UsersHTTPHandler) ProtectedRoutes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/users/me",
			Handler:    h.GetCurrentUser,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/users/logout",
			Handler:    h.Logout,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/users",
			Handler:    h.GetUsers,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/users/{id}",
			Handler:    h.GetUser,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/users/{id}",
			Handler:    h.DeleteUser,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/users/{id}",
			Handler:    h.PatchUser,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
	}
}
