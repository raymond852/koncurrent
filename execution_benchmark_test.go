package koncurrent

import (
	"context"
	"testing"
)

func BenchmarkExecuteSerial(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1 = NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		var task2 = NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		var task3= NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		_, _ = ExecuteSerial(task1, task2, task3).Await(context.Background())
	}
}

func BenchmarkExecuteParallel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1 = NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		var task2 = NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		var task3= NewTask(func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}, NewAsyncExecutor())

		_, _ = ExecuteParallel(task1, task2, task3).Await(context.Background())
	}
}
