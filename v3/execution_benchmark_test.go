package koncurrent

import (
	"context"
	"testing"
)

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

		_, _ = ExecuteSerial(
			task1Func.Immediate().Recover(),
			task2Func.Immediate().Recover(),
			task3Func.Immediate().Recover()).
			Await(context.Background())
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

		_, _ = ExecuteSerial(
			task1Func.Async().Recover(),
			task2Func.Async().Recover(),
			task3Func.Async().Recover()).
			Await(context.Background())
	}
}

func BenchmarkExecuteSerial_Pool(b *testing.B) {
	var pe = NewPoolExecutor(2, 10)
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

		_, _ = ExecuteSerial(
			task1Func.Pool(pe).Recover(),
			task2Func.Pool(pe).Recover(),
			task3Func.Pool(pe).Recover()).
			Await(context.Background())
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

		_, _ = ExecuteParallel(
			task1Func.Immediate().Recover(),
			task2Func.Immediate().Recover(),
			task3Func.Immediate().Recover()).
			Await(context.Background())
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

		_, _ = ExecuteParallel(
			task1Func.Async().Recover(),
			task2Func.Async().Recover(),
			task3Func.Async().Recover()).
			Await(context.Background())
	}
}

func BenchmarkExecuteParallel_Pool(b *testing.B) {
	var pe = NewPoolExecutor(2, 10)
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

		_, _ = ExecuteParallel(
			task1Func.Pool(pe).Recover(),
			task2Func.Pool(pe).Recover(),
			task3Func.Pool(pe).Recover()).
			Await(context.Background())
	}
}