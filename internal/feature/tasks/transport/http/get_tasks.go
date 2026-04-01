package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_utils "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, userID, err := getTasksQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get query params",
		)
		return
	}
	tasksDomain, err := h.tasksService.GetTasks(
		ctx,
		limit,
		offset,
		userID,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(TasksDTOFromDomains(tasksDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getTasksQueryParams(r *http.Request) (*int, *int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`parameter "limit": %w`, err)
	}
	offset, err := core_http_utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`parameter "offset": %w`, err)
	}
	userID, err := core_http_utils.GetIntQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`parameter "user_id": %w`, err)
	}
	return limit, offset, userID, nil
}
