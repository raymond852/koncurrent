package koncurrent

import (
	"context"
)

type ImmediateExecutor struct {
}

type immediateFuture struct {
	err error
}

func (p *immediateFuture) Get() error {
	return p.err
}

func NewImmediateExecutor() TaskExecutor {
	return &ImmediateExecutor{}
}

func (a *ImmediateExecutor) Execute(ctx context.Context, task TaskFunc) TaskFuture {
	return &immediateFuture{err: task(ctx)}
}
