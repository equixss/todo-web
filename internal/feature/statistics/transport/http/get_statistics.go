package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	core_http_request "github.com/equixss/todo-web/internal/core/transport/http/request"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated              int      `json:"tasks_created"`
	TasksCompleted            int      `json:"tasks_completed"`
	TasksCompletedPercent     *float64 `json:"tasks_completed_percent"`
	TasksAverageCompletedTime *string  `json:"tasks_average_completed_time"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	queryParams, err := getQueryParameters(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}
	statistics, err := h.statisticsService.GetStatistics(
		ctx,
		queryParams.UserID,
		queryParams.From,
		queryParams.To,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}
	responseHandler.JSONResponse(domainStatisticsToDTO(statistics), http.StatusOK)
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
	UserID *int
	From   *time.Time
	To     *time.Time
}

func getQueryParameters(r *http.Request) (*queryParams, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'to' query param: %w", err)
	}
	return &queryParams{
		UserID: userID,
		From:   from,
		To:     to,
	}, nil
}
