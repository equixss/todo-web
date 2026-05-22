package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type GetTasksResponse []TaskDTOResponse

// @Summary Получение списка задач пользователя
// @Description Возвращает список всех задач текущего пользователя. Требуется авторизация.
// @Tags tasks
// @Accept json
// @Produce json
// @Param limit query int false "Лимит количества записей"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {object} GetTasksResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(c *gin.Context) {

	userID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	limit, offset, err := getTasksQueryParams(c.Request)
	if err != nil {
		h.presenter.ErrorResponse(
			c,
			err,
			"failed to get query params",
		)
		return
	}

	tasksDomain, err := h.tasksService.GetTasks(
		c.Request.Context(),
		limit,
		offset,
		userID,
	)
	if err != nil {
		h.presenter.ErrorResponse(
			c,
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(TasksDTOFromDomains(tasksDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}

func getTasksQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf(`parameter "limit": %w`, err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf(`parameter "offset": %w`, err)
	}
	return limit, offset, nil
}
