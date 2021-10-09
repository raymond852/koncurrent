package koncurrent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"runtime/debug"
)

type ImmediateExecutor struct {
}

func (p ImmediateExecutor) Execute(ctx context.Context, taskFunc TaskFunc, taskId int, resultChn chan TaskResult, opt TaskExecutionOptions) {
	defer func() {
		if r := recover(); r != nil {
			resultChn <- TaskResult{
				err : PanicError{
					Stack: debug.Stack(),
				},
				id: taskId,
			}
		}
	}()
	c := ctx
	if len(opt.tracingSpanName) > 0 {
		span, spanCtx := opentracing.StartSpanFromContext(ctx, opt.tracingSpanName)
		c = spanCtx
		defer span.Finish()
	}
	resultChn <- TaskResult{
		err: taskFunc(c),
		id:  taskId,
	}

}