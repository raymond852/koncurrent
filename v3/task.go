package koncurrent

import (
	"context"
)

const (
	taskExecutionTypeImmediate = iota
	taskExecutionTypeNewGoRoutine
	taskExecutionTypeOther
)

type TaskFunc func(ctx context.Context) error

type TaskFuture chan error

type TaskExecutor interface {
	Execute(ctx context.Context, task TaskFunc) TaskFuture
}

type TaskExecution struct {
	executionType int
	taskFunc      TaskFunc
	executor      TaskExecutor
}

func (t TaskFunc) Immediate() TaskExecution {
	return TaskExecution{
		executionType: taskExecutionTypeImmediate,
		taskFunc:      t,
	}
}

func (t TaskFunc) Async() TaskExecution {
	return TaskExecution{
		executionType: taskExecutionTypeNewGoRoutine,
		taskFunc:      t,
	}
}

func (t TaskFunc) Pool(executor PoolExecutor) TaskExecution {
	return TaskExecution{
		executionType: taskExecutionTypeOther,
		taskFunc:      t,
		executor:      executor,
	}
}
