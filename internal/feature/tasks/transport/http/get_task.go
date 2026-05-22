package tasks_transport_http

import (
	"errors"
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type GetTaskResponse TaskDTOResponse

// @Summary Получение задачи по ID
// @Description Возвращает данные задачи по её ID. Требуется авторизация.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(c *gin.Context) {

	userID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	taskID, err := core_http_request.GetIntPathValue(c.Request, "id")
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get ID path param")
		return
	}

	taskDomain, err := h.tasksService.GetTask(c.Request.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			h.presenter.ErrorResponse(c, err, "task not found")
			return
		}
		h.presenter.ErrorResponse(c, err, "failed to get task")
		return
	}

	response := GetTaskResponse(TaskDTOFromDomain(taskDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}
