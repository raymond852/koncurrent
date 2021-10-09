package koncurrent

import (
	"context"
	"testing"
)

var pe = NewPoolExecutor(2, 10)

func BenchmarkExecuteSerial_Immediate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteSerial(task1Func.Immediate(), task2Func.Immediate(), task3Func.Immediate()).Await(context.Background())
	}
}

func BenchmarkExecuteSerial_Async(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteSerial(task1Func.Async(), task2Func.Async(), task3Func.Async()).Await(context.Background())
	}
}

func BenchmarkExecuteSerial_Pool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteSerial(task1Func.Pool(pe), task2Func.Pool(pe), task3Func.Pool(pe)).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Immediate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteParallel(task1Func.Immediate(), task2Func.Immediate(), task3Func.Immediate()).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Async(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteParallel(task1Func.Async(), task2Func.Async(), task3Func.Async()).Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Pool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var task1Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task2Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		var task3Func TaskFunc = func(ctx context.Context) error {
			return nil
		}

		_, _ = ExecuteParallel(task1Func.Pool(pe), task2Func.Pool(pe), task3Func.Pool(pe)).Await(context.Background())
	}
}