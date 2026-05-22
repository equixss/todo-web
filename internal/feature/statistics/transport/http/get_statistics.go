package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	"github.com/gin-gonic/gin"
)

type GetStatisticsResponse struct {
	TasksCreated              int      `json:"tasks_created"`
	TasksCompleted            int      `json:"tasks_completed"`
	TasksCompletedPercent     *float64 `json:"tasks_completed_percent"`
	TasksAverageCompletedTime *string  `json:"tasks_average_completed_time"`
}

// @Summary Получение статистики задач пользователя
// @Description Возвращает статистику выполнения задач текущего пользователя. Требуется авторизация.
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} GetStatisticsResponse
// @Failure 401 {object} map[string]string "Требуется авторизация"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(c *gin.Context) {

	userID, ok := core_http_middleware.GetUserIDFromContext(c.Request.Context())
	if !ok {
		h.presenter.ErrorResponse(c, core_errors.ErrUnauthorized, "authentication required")
		return
	}

	queryParams, err := getQueryParameters(c.Request)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get query params")
		return
	}
	statistics, err := h.statisticsService.GetStatistics(
		c.Request.Context(),
		userID,
		queryParams.From,
		queryParams.To,
	)
	if err != nil {
		h.presenter.ErrorResponse(c, err, "failed to get statistics")
		return
	}
	h.presenter.JSONResponse(c, domainStatisticsToDTO(statistics), http.StatusOK)
}

func domainStatisticsToDTO(statistics domain.Statistics) GetStatisticsResponse {
	var tasksAvgTime *string
	if statistics.TasksAverageCompletedTime != nil {
		timeStr := statistics.TasksAverageCompletedTime.String()
		tasksAvgTime = &timeStr
	}
	return GetStatisticsResponse{
		TasksCreated:              statistics.TasksCreated,
		TasksCompleted:            statistics.TasksCompleted,
		TasksCompletedPercent:     statistics.TasksCompletedPercent,
		TasksAverageCompletedTime: tasksAvgTime,
	}
}

type queryParams struct {
	From *time.Time
	To   *time.Time
}

func getQueryParameters(r *http.Request) (*queryParams, error) {
	const (
		fromQueryParamKey = "from"
		toQueryParamKey   = "to"
	)
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'to' query param: %w", err)
	}
	return &queryParams{
		From: from,
		To:   to,
	}, nil
}
