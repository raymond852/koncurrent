package koncurrent

import (
	"context"
	"testing"
)

func TestAsyncExecutor_Execute(t *testing.T) {
	underTest := NewAsyncExecutor()
	ret := underTest.Execute(context.Background(), func(ctx context.Context) *TaskResult {
		return NewTaskResult("test", nil)
	})
	result := ret.Get()
	assertNil(t, result.Error)
	assertEqual(t, "test", result.Result)
}
