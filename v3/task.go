package koncurrent

import (
	"context"
)

type TaskFunc func(ctx context.Context) error

type PanicError struct {
	Stack []byte
}

func (e PanicError) Error() string {
	return "panic:" + string(e.Stack)
}

type TaskExecutor interface {
	Execute(ctx context.Context, taskFunc TaskFunc, taskId int, resultChn chan TaskResult, taskExecutionOpts TaskExecutionOptions)
}

type TaskResult struct {
	err error
	id  int
}

type TaskExecution struct {
	options  TaskExecutionOptions
	taskFunc TaskFunc
	executor TaskExecutor
}

type TaskExecutionOptions struct {
	tracingSpanName  string
	recoverFromPanic bool
}

func (t TaskExecution) Recover() TaskExecution {
	ret := t
	ret.options.recoverFromPanic = true
	return ret
}

func (t TaskExecution) Tracing(spanName string) TaskExecution {
	ret := t
	ret.options.tracingSpanName = spanName
	return ret
}

func (t TaskExecution) Execution() Execution {
	return ExecuteSerial(t)
}

func (t TaskFunc) Immediate() TaskExecution {
	return TaskExecution{
		taskFunc: t,
		executor: ImmediateExecutor{},
	}
}

func (t TaskFunc) Async() TaskExecution {
	return TaskExecution{
		taskFunc: t,
		executor: AsyncExecutor{},
	}
}

func (t TaskFunc) Pool(executor PoolExecutor) TaskExecution {
	return TaskExecution{
		taskFunc: t,
		executor: executor,
	}
}
