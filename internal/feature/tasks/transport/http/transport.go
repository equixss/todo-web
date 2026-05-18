package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_response_presenter "github.com/equixss/todo-web/internal/core/transport/http/response"
	core_http_server "github.com/equixss/todo-web/internal/core/transport/http/server"
)

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
	GetTask(
		ctx context.Context,
		id int,
		userID int,
	) (domain.Task, error)
	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
		userID int,
	) (domain.Task, error)
	DeleteTask(
		ctx context.Context,
		id int,
		userID int,
	) error
	GetTasks(
		ctx context.Context,
		limit *int,
		offset *int,
		userID int,
	) ([]domain.Task, error)
}

type TasksHTTPHandler struct {
	presenter    *core_http_response_presenter.HTTPResponsePresenter
	tasksService TasksService
	jwtMW        *core_http_middleware.JWTMiddleware
}

func NewTasksHttpHandler(
	tasksService TasksService,
	jwtMW *core_http_middleware.JWTMiddleware,
	presenter *core_http_response_presenter.HTTPResponsePresenter,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: tasksService, jwtMW: jwtMW, presenter: presenter}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/tasks",
			Handler:    h.CreateTask,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks",
			Handler:    h.GetTasks,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks/{id}",
			Handler:    h.GetTask,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/tasks/{id}",
			Handler:    h.DeleteTask,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/tasks/{id}",
			Handler:    h.PatchTask,
			Middleware: []core_http_middleware.Middleware{h.jwtMW.Authenticate()},
		},
	}
}
