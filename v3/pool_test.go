package koncurrent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func TestPoolExecutor_Execute(t *testing.T) {
	underTest := NewPoolExecutor(10, 10)
	var span opentracing.Span
	resultChan := make(chan TaskResult)
	underTest.Execute(context.Background(), func(ctx context.Context) error {
		span = opentracing.SpanFromContext(ctx)
		panic("test panic")
	}, 0, resultChan, TaskExecutionOptions{
		tracingSpanName:  "test",
		recoverFromPanic: true,
	})
	result := <-resultChan
	_, ok := result.err.(PanicError)
	assertTrue(t, ok)
	assertTrue(t, span != nil)
	assertTrue(t, len(result.err.Error()) > 0)
}
