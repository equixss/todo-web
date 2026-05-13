package tasks_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=100"`
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.GetUserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(ErrUnauthorized, "authentication required")
		return
	}

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		userID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	response := CreateTaskResponse(TaskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

var ErrUnauthorized = core_errors.ErrUnauthorized
