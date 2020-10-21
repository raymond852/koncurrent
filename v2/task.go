package koncurrent

import (
	"context"
)

type TaskFunc func(ctx context.Context) error

type TaskFuture interface {
	Get() error
}

type TaskExecutor interface {
	Execute(ctx context.Context, task TaskFunc) TaskFuture
}

type Task struct {
	taskFunc TaskFunc
	executor TaskExecutor
}

func NewTask(task TaskFunc, executor TaskExecutor) *Task {
	return &Task{
		taskFunc: task,
		executor: executor,
	}
}
