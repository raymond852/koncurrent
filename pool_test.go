package koncurrent

import (
	"context"
	"testing"
)

func TestPoolExecutor_Execute(t *testing.T) {
	underTest := NewPoolExecutor(10, 10)
	ret := underTest.Execute(context.Background(), func(ctx context.Context) *TaskResult {
		return NewTaskResult("test", nil)
	})
	result := ret.Get()
	assertNil(t, result.Error)
	assertEqual(t, "test", result.Result)
}
