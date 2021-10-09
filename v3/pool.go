package koncurrent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"runtime/debug"
)

type PoolExecutor struct {
	queue chan taskContext
}

type taskContext struct {
	context.Context
	opt       TaskExecutionOptions
	task      TaskFunc
	taskId    int
	resultChn chan TaskResult
}

func (p PoolExecutor) Execute(ctx context.Context, task TaskFunc, taskId int, resultChn chan TaskResult, opt TaskExecutionOptions) {
	p.queue <- taskContext{
		Context:   ctx,
		task:      task,
		taskId:    taskId,
		resultChn: resultChn,
		opt:       opt,
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
				ctx := taskCtx.Context
				tracingSpanName := taskCtx.opt.tracingSpanName
				resultChn := taskCtx.resultChn
				taskFunc := taskCtx.task
				taskId := taskCtx.taskId
				func() {
					defer func() {
						if r := recover(); r != nil {
							resultChn <- TaskResult{
								err: PanicError{
									Stack: debug.Stack(),
								},
								id: taskId,
							}
						}
					}()
					var s opentracing.Span
					c := ctx
					if len(tracingSpanName) > 0 {
						span, spanCtx := opentracing.StartSpanFromContext(ctx, tracingSpanName)
						s = span
						c = spanCtx
					}
					taskErr := taskFunc(c)
					resultChn <- TaskResult{
						err: taskErr,
						id:  taskId,
					}
					if s != nil {
						s.Finish()
					}
				}()
			}
		}()
	}
	return ret
}
