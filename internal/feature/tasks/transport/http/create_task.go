package tasks_transport_http

import (
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=100"`
}

type CreateTaskResponse TaskDTOResponse

// @Summary Создание новой задачи
// @Description Создает новую задачу. Требуется авторизация.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Запрос на создание задачи"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} map[string]string "Некорректные данные запроса"
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(c *gin.Context) {

	userID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate HTTP request")
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		userID,
	)

	taskDomain, err := h.tasksService.CreateTask(c.Request.Context(), taskDomain)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to create task")
		return
	}

	response := CreateTaskResponse(TaskDTOFromDomain(taskDomain))
	h.presenter.JSONResponse(c, response, http.StatusCreated)
}
