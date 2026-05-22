package tasks_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_types "github.com/equixss/todo-web/internal/core/transport/http/types"
	"github.com/gin-gonic/gin"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type PatchTaskSwaggerRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *string `json:"completed"`
}

type PatchTaskResponse TaskDTOResponse

// @Summary Обновление задачи
// @Description Частичное обновление данных задачи. Требуется авторизация.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body PatchTaskSwaggerRequest true "Запрос на обновление"
// @Success 200 {object} PatchTaskResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(c *gin.Context) {
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

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(c.Request, &request); err != nil {
		h.presenter.ErrorResponse(c, err, "failed to decode and validate")
		return
	}
	taskPatch := taskPatchFromRequest(request)
	taskDomain, err := h.tasksService.PatchTask(c.Request.Context(), taskID, taskPatch, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			h.presenter.ErrorResponse(c, err, "task not found")
			return
		}
		h.presenter.ErrorResponse(c, err, "failed to patch task")
		return
	}
	response := PatchTaskResponse(TaskDTOFromDomain(taskDomain))
	h.presenter.JSONResponse(c, response, http.StatusOK)
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("field 'title' cannot be nil")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("field 'title' length must be between 1 and 100 symbols")
		}
	}
	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 100 {
				return fmt.Errorf("field 'description' length must be between 1 and 100 symbols")
			}
		}
	}
	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("field 'completed' cannot be nil")
		}
	}
	return nil
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
