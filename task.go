package koncurrent

import (
	"context"
)

type TaskFunc func(ctx context.Context) *TaskResult

type TaskResult struct {
	Result interface{}
	Error  error
}

func NewTaskResult(result interface{}, err error) *TaskResult {
	return &TaskResult{
		Result: result,
		Error:  err,
	}
}

type TaskFuture interface {
	Get() *TaskResult
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
