package koncurrent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestExecuteSerial(t *testing.T) {
	var time1, time2 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor(), NewImmediateExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i])).Await(context.Background())
		results := iter.Next()
		assertNil(t, err)
		assertEqual(t, 2, len(results))
		assertNotNil(t, time1)
		assertNotNil(t, time2)
		assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
		assertEqual(t, "task1", results[0].Result)
		assertEqual(t, "task2", results[1].Result)
	}
}

func TestExecuteSerial_CanAbortWithError(t *testing.T) {
	var time1, time3 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		return NewTaskResult(nil, errors.New( "test"))
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor(), NewImmediateExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i]), NewTask(t3, executors[i])).Await(context.Background())
		results := iter.Next()
		assertNotNil(t, err)
		assertEqual(t, 2, len(results))
		assertNotNil(t, time1)
		assertNil(t, time3)
		assertEqual(t, "task1", results[0].Result)
		assertNotNil(t, results[1].Error)
	}
}

func TestExecuteParallel(t *testing.T) {
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(5000 * time.Millisecond)
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(1000 * time.Millisecond)
		return NewTaskResult("task2", nil)
	}

	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		now := time.Now()
		iter, err := ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i])).Await(context.Background())
		elapse := time.Now().Sub(now)
		assertTrue(t, elapse < 5100*time.Millisecond)
		results := iter.Next()
		assertNil(t, err)
		assertEqual(t, 2, len(results))
		assertEqual(t, "task1", results[0].Result)
		assertEqual(t, "task2", results[1].Result)
	}
}

func TestExecuteParallel_ImmediateExecutor(t *testing.T) {
	var time1, time2 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}

	iter, err := ExecuteParallel(NewTask(t1, NewImmediateExecutor()), NewTask(t2, NewImmediateExecutor())).Await(context.Background())
	results := iter.Next()
	assertNil(t, err)
	assertEqual(t, 2, len(results))
	assertNotNil(t, time1)
	assertNotNil(t, time2)
	assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
	assertEqual(t, "task1", results[0].Result)
	assertEqual(t, "task2", results[1].Result)
}

func TestExecuteParallel_Error(t *testing.T) {
	var time1, time3 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		return NewTaskResult(nil, errors.New("test"))
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task2", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		return NewTaskResult(nil, errors.New("test"))
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i]), NewTask(t3, executors[i]), NewTask(t4, executors[i])).Await(context.Background())
		assertNotNil(t, iter)
		results := iter.Next()
		fmt.Println(fmt.Errorf("err %w", err))
		assertNotNil(t, err)
		assertEqual(t, 4, len(results))
		assertNotNil(t, time1)
		assertNotNil(t, time3)
		assertEqual(t, "task1", results[0].Result)
		assertEqual(t, "task2", results[2].Result)
		assertNotNil(t, results[1].Error)
		assertNotNil(t, results[3].Error)
	}
}

func TestCascade_SerialAfterParallel(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return NewTaskResult("task4", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteParallel(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			Await(context.Background())
		assertNil(t, err)
		i := 0
		for {
			results := iter.Next()
			if results == nil {
				break
			}
			assertEqual(t, 2, len(results))
			assertNotNil(t, time1)
			assertNotNil(t, time2)
			assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
			assertNotNil(t, time3)
			assertTrue(t, (*time3).Sub(*time2) >= 100*time.Millisecond)
			assertNotNil(t, time4)
			assertTrue(t, (*time4).Sub(*time3) < 10*time.Millisecond)
			if i == 0 {
				assertEqual(t, "task1", results[0].Result)
				assertEqual(t, "task2", results[1].Result)
			} else if i == 1 {
				assertEqual(t, "task3", results[0].Result)
				assertEqual(t, "task4", results[1].Result)
			}
			i++
		}
		assertEqual(t, 2, i)
	}
}

func TestCascade_ParallelAfterSerial(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return NewTaskResult("task4", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteSerial(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			Await(context.Background())
		assertNil(t, err)
		i := 0
		for {
			results := iter.Next()
			if results == nil {
				break
			}
			assertEqual(t, 2, len(results))
			assertNotNil(t, time1)
			assertNotNil(t, time2)
			assertTrue(t, (*time2).Sub(*time1) <= 10*time.Millisecond)
			assertNotNil(t, time3)
			assertTrue(t, (*time3).Sub(*time2) >= 100*time.Millisecond)
			assertNotNil(t, time4)
			assertTrue(t, (*time4).Sub(*time3) >= 100*time.Millisecond)
			if i == 0 {
				assertEqual(t, "task1", results[0].Result)
				assertEqual(t, "task2", results[1].Result)
			} else if i == 1 {
				assertEqual(t, "task3", results[0].Result)
				assertEqual(t, "task4", results[1].Result)
			}
			i++
		}
		assertEqual(t, 2, i)
	}
}

func TestCascade_SerialAfterParallelError(t *testing.T) {
	var time1, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		return NewTaskResult(nil, errors.New("test"))
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return NewTaskResult("task4", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteParallel(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			Await(context.Background())
		assertNotNil(t, err)
		i := 0
		for {
			results := iter.Next()
			if results == nil {
				break
			}
			assertEqual(t, 2, len(results))
			assertNotNil(t, time1)
			assertNil(t, time3)
			assertNil(t, time4)
			assertEqual(t, "task1", results[0].Result)
			assertNil(t, results[1].Result)
			assertNotNil(t, results[1].Error)
			i++
		}
		assertEqual(t, 1, i)
	}
}

func TestCascade_SerialAfterParallelAfterParallel(t *testing.T) {
	var time1, time2, time3, time4, time5, time6 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return NewTaskResult("task4", nil)
	}
	var t5 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time5 = &now
		return NewTaskResult("task5", nil)
	}
	var t6 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time6 = &now
		return NewTaskResult("task6", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteParallel(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			ExecuteParallel(NewTask(t5, executors[i]), NewTask(t6, executors[i])).
			Await(context.Background())
		assertNil(t, err)
		i := 0
		for {
			results := iter.Next()
			if results == nil {
				break
			}
			assertEqual(t, 2, len(results))
			assertNotNil(t, time1)
			assertNotNil(t, time2)
			assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
			assertNotNil(t, time3)
			assertTrue(t, (*time3).Sub(*time2) >= 100*time.Millisecond)
			assertNotNil(t, time4)
			assertTrue(t, (*time4).Sub(*time3) < 10*time.Millisecond)
			assertNotNil(t, time5)
			assertTrue(t, (*time5).Sub(*time4) >= 100*time.Millisecond)
			assertNotNil(t, time6)
			assertTrue(t, (*time6).Sub(*time5) < 10*time.Millisecond)
			if i == 0 {
				assertEqual(t, "task1", results[0].Result)
				assertEqual(t, "task2", results[1].Result)
			} else if i == 1 {
				assertEqual(t, "task3", results[0].Result)
				assertEqual(t, "task4", results[1].Result)
			} else if i == 2 {
				assertEqual(t, "task5", results[0].Result)
				assertEqual(t, "task6", results[1].Result)
			}
			i++
		}
		assertEqual(t, 3, i)
	}
}

func TestExecution_Async(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return NewTaskResult("task1", nil)
	}
	var t2 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return NewTaskResult("task2", nil)
	}
	var t3 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return NewTaskResult("task3", nil)
	}
	var t4 TaskFunc = func(ctx context.Context) *TaskResult {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return NewTaskResult("task4", nil)
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		wg := sync.WaitGroup{}
		wg.Add(1)
		ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteSerial(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			Async(context.Background(), executors[i], func(iter *ExecutionResultIterator, err error) {
				assertNil(t, err)
				i := 0
				for {
					results := iter.Next()
					if results == nil {
						break
					}
					assertEqual(t, 2, len(results))
					assertNotNil(t, time1)
					assertNotNil(t, time2)
					assertTrue(t, (*time2).Sub(*time1) <= 10*time.Millisecond)
					assertNotNil(t, time3)
					assertTrue(t, (*time3).Sub(*time2) >= 100*time.Millisecond)
					assertNotNil(t, time4)
					assertTrue(t, (*time4).Sub(*time3) >= 100*time.Millisecond)
					if i == 0 {
						assertEqual(t, "task1", results[0].Result)
						assertEqual(t, "task2", results[1].Result)
					} else if i == 1 {
						assertEqual(t, "task3", results[0].Result)
						assertEqual(t, "task4", results[1].Result)
					}
					i++
				}
				assertEqual(t, 2, i)
				wg.Done()
			})
		wg.Wait()
	}
}
