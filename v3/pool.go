package koncurrent

import (
	"context"
)

type PoolExecutor struct {
	queue chan taskContext
}

type taskContext struct {
	context.Context
	task      TaskFunc
	taskId    int
	resultChn chan TaskResult
}

func (p PoolExecutor) Execute(ctx context.Context, task TaskFunc, taskId int, resultChn chan TaskResult) {
	p.queue <- taskContext{
		Context:   ctx,
		task:      task,
		taskId:    taskId,
		resultChn: resultChn,
	}
}

func NewPoolExecutor(poolSize int, queueSize int) PoolExecutor {
	ret := PoolExecutor{
		queue: make(chan taskContext, queueSize),
	}
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				taskCtx := <-ret.queue
				taskErr := taskCtx.task(taskCtx.Context)
				taskCtx.resultChn <- TaskResult{
					err: taskErr,
					id:  taskCtx.taskId,
				}
			}
		}()
	}
	return ret
}
