package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_server "github.com/equixss/todo-web/internal/core/transport/http/server"
)

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

type TasksHTTPHandler struct {
	tasksService TasksService
}

func NewTasksHttpHandler(tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: tasksService}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
	}
}
