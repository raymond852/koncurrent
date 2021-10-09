package koncurrent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func TestImmediateExecutor_Execute(t *testing.T) {
	resultChan := make(chan TaskResult, 1)
	var span opentracing.Span
	ImmediateExecutor{}.Execute(context.Background(), func(ctx context.Context) error {
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
