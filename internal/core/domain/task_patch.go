package domain

import (
	"fmt"
	"time"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	tmp := *t
	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid patch task: %w", err)
	}

	*t = tmp
	return nil
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("'Title' can't be patched to null:%w", core_errors.ErrInvalidArgument)
	}
	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("'Completed' can't be patched to null:%w", core_errors.ErrInvalidArgument)
	}
	return nil
}
