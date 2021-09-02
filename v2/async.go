package koncurrent

import "context"

var instance AsyncExecutor

type AsyncExecutor struct {
}

type channelTaskFuture struct {
	chn chan error
}

func (c channelTaskFuture) Get() error {
	return <-c.chn
}

func (a AsyncExecutor) Execute(ctx context.Context, task TaskFunc) TaskFuture {
	chn := make(chan error, 1)
	go func() {
		chn <- task(ctx)
		close(chn)
	}()
	return channelTaskFuture{chn: chn}
}

func NewAsyncExecutor() TaskExecutor {
	return instance
}
