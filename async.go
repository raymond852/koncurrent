package koncurrent

import "context"

var instance *AsyncExecutor

type AsyncExecutor struct {
}

func init() {
	instance = &AsyncExecutor{}
}

type channelTaskFuture struct {
	chn chan *TaskResult
}

func (c *channelTaskFuture) Get() *TaskResult {
	return <-c.chn
}

func (a *AsyncExecutor) Execute(ctx context.Context, task TaskFunc) TaskFuture {
	chn := make(chan *TaskResult, 1)
	go func() {
		chn <- task(ctx)
		close(chn)
	}()
	return &channelTaskFuture{chn: chn}
}

func NewAsyncExecutor() TaskExecutor {
	return instance
}
