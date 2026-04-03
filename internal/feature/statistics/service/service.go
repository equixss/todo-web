package statistics_service

import (
	"context"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
)

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

func NewStatisticsService(repository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		repository,
	}
}
