package domain

import "time"

type Statistics struct {
	TasksCreated              int
	TasksCompleted            int
	TasksCompletedPercent     *float64
	TasksAverageCompletedTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompletedPercent *float64,
	tasksAverageCompletedTime *time.Duration,
) Statistics {
	return Statistics{
		tasksCreated,
		tasksCompleted,
		tasksCompletedPercent,
		tasksAverageCompletedTime,
	}
}
