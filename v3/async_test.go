package koncurrent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func TestAsyncExecutor_Execute(t *testing.T) {
	resultChan := make(chan TaskResult)
	var span opentracing.Span
	AsyncExecutor{}.Execute(context.Background(), func(ctx context.Context) error {
		span = opentracing.SpanFromContext(ctx)
		panic("test panic")
	}, 1, resultChan, TaskExecutionOptions{
		tracingSpanName:  "test",
		recoverFromPanic: true,
	})
	result := <-resultChan
	_, ok := result.err.(PanicError)
	assertTrue(t, ok)
	assertTrue(t, span != nil)
	assertTrue(t, len(result.err.Error()) > 0)
}
