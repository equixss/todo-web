package domain

import (
	"fmt"
	"time"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

type Task struct {
	ID      int `json:"id"`
	Version int `json:"version"`

	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`

	AuthorUserID int `json:"author_user_id"`
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}
	if t.CompletedAt == nil {
		return nil
	}
	duration := t.CompletedAt.Sub(t.CreatedAt)
	return &duration
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid 'Title' len %d:%w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}
	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 100 {
			return fmt.Errorf(
				"invalid 'Description' len %d:%w",
				titleLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}
	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"CompletedAt can't be nil if Completed is true:%w",
				core_errors.ErrInvalidArgument,
			)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"CompletedAt can't be before CreatedAt:%w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"CompletedAt must be nil if Completed is false:%w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil
}
