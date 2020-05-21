package koncurrent

import (
	"context"
)

type ImmediateExecutor struct {
}

type immediateFuture struct {
	result *TaskResult
}

func (p *immediateFuture) Get() *TaskResult {
	return p.result
}

func NewImmediateExecutor() TaskExecutor {
	return &ImmediateExecutor{}
}

func (a *ImmediateExecutor) Execute(ctx context.Context, task TaskFunc) TaskFuture {
	result := task(ctx)
	return &immediateFuture{result: result}
}
