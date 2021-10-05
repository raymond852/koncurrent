package koncurrent

import (
	"context"
	"testing"
)

var pe = NewPoolExecutor(2, 10)

func BenchmarkExecuteSerial_Immediate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteSerial(NewTask(task1Func, NewImmediateExecutor()), NewTask(task2Func, NewImmediateExecutor()), NewTask(task3Func, NewImmediateExecutor())).Await(context.Background())
	}
}

func BenchmarkExecuteSerial_Async(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteSerial(NewTask(task1Func, NewAsyncExecutor()), NewTask(task2Func, NewAsyncExecutor()), NewTask(task3Func, NewAsyncExecutor())).Await(context.Background())
	}
}

func BenchmarkExecuteSerial_Pool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteSerial(NewTask(task1Func, pe), NewTask(task2Func, pe), NewTask(task3Func, pe)).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Immediate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteSerial(NewTask(task1Func, NewImmediateExecutor()), NewTask(task2Func, NewImmediateExecutor()), NewTask(task3Func, NewImmediateExecutor())).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Async(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteParallel(NewTask(task1Func, NewAsyncExecutor()), NewTask(task2Func, NewAsyncExecutor()), NewTask(task3Func, NewAsyncExecutor())).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Pool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task2Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		var task3Func TaskFunc = func(ctx context.Context) *TaskResult {
			return NewTaskResult(nil, nil)
		}

		_, _ = ExecuteSerial(NewTask(task1Func, pe), NewTask(task2Func, pe), NewTask(task3Func, pe)).Await(context.Background())
	}
}