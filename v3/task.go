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

type TaskResult struct {
	err error
	id  int
}

type TaskExecution struct {
	executionType int
	taskFunc      TaskFunc
	executor      PoolExecutor
}

func (t TaskExecution) Execution() Execution {
	return ExecuteSerial(t)
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
