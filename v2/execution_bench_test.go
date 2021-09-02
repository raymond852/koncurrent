package koncurrent

import (
	"context"
	"testing"
)

func TestNewTask(b *testing.T) {
	//for i := 0; i < b.N; i++ {
		var task1Result string
		var task2Result string
		var task3Result string
		var task1 = NewTask(func(ctx context.Context) error {
			task1Result = "task1"
			return nil
		}, NewAsyncExecutor())

		var task2 = NewTask(func(ctx context.Context) error {
			task2Result = "task2"
			return nil
		}, NewAsyncExecutor())

		var task3 = NewTask(func(ctx context.Context) error {
			task3Result = "task3"
			return nil
		}, NewAsyncExecutor())

		println(task1Result)
		println(task2Result)
		println(task3Result)

		_, _ = ExecuteSerial(task1, task2, task3).Await(context.Background())
	//}
}
