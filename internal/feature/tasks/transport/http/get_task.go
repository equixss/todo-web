package tasks_transport_http

import (
	"errors"
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_utils "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.GetUserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(ErrUnauthorized, "authentication required")
		return
	}

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get ID path param")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			responseHandler.ErrorResponse(err, "task not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	response := GetTaskResponse(TaskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
