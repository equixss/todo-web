package tasks_transport_http

import (
	"errors"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type DeleteTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) DeleteTask(c *gin.Context) {

	authUserID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	taskID, err := core_http_request.GetIntPathValue(c.Request, "id")
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get ID path param")
		return
	}

	if err := h.tasksService.DeleteTask(c.Request.Context(), taskID, authUserID); err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			h.presenter.ErrorResponse(c, err, "task not found")
			return
		}
		h.presenter.ErrorResponse(c, err, "failed to delete task")
		return
	}
	h.presenter.ResponseNoContent(c)
}
