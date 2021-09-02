package koncurrent

import (
	"context"
)

type PoolExecutor struct {
	queue chan taskContext
}

type taskContext struct {
	context.Context
	task   TaskFunc
	result chan error
}

func (p PoolExecutor) Execute(ctx context.Context, task TaskFunc) TaskFuture {
	chn := make(chan error, 1)
	p.queue <- taskContext{
		Context: ctx,
		task:    task,
		result:  chn,
	}
	return channelTaskFuture{chn: chn}
}

func NewPoolExecutor(poolSize int, queueSize int) TaskExecutor {
	ret := PoolExecutor{
		queue: make(chan taskContext, queueSize),
	}
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				taskCtx := <-ret.queue
				result := taskCtx.task(taskCtx.Context)
				taskCtx.result <- result
				close(taskCtx.result)
			}
		}()
	}
	return ret
}
