package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{},
				fmt.Errorf(
					"from '%v' must be before to '%v':%w",
					*from,
					*to,
					core_errors.ErrInvalidArgument,
				)
		}
	}
	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{},
			fmt.Errorf("get tasks from repository:%w", err)
	}

	statistics := calculateStatistic(tasks)

	return statistics, nil
}

func calculateStatistic(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}
	tasksCreated := len(tasks)
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
			completionDuration := task.CompletionDuration()
			if completionDuration != nil {
				totalCompletionDuration += *completionDuration
			}
		}
	}
	tasksCompletedPercent := (float64(tasksCompleted) / float64(tasksCreated)) * 100.0

	var avgTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		avgTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedPercent,
		avgTime,
	)
}
