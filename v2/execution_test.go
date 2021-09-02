package koncurrent

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestExecuteSerial(t *testing.T) {
	var time1, time2 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
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
	}
}

func TestExecuteSerial_CanAbortWithError(t *testing.T) {
	var time1, time3 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		return errors.New("test")
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor(), NewImmediateExecutor()}
	for i := range executors {
		iter, err := ExecuteSerial(NewTask(t1, executors[i]), NewTask(t2, executors[i]), NewTask(t3, executors[i])).Await(context.Background())
		resultErrors := iter.Next()
		assertNotNil(t, err)
		assertEqual(t, 2, len(resultErrors))
		assertNotNil(t, time1)
		assertNil(t, time3)
		assertNotNil(t, resultErrors[1])
	}
}

func TestExecuteParallel(t *testing.T) {
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(5000 * time.Millisecond)
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		now := time.Now()
		iter, err := ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i])).Await(context.Background())
		elapse := time.Now().Sub(now)
		assertTrue(t, elapse < 5100*time.Millisecond)
		resultErrors := iter.Next()
		assertNil(t, err)
		assertEqual(t, 2, len(resultErrors))
		assertNil(t, resultErrors[0])
		assertNil(t, resultErrors[1])
	}
}

func TestExecuteParallel_ImmediateExecutor(t *testing.T) {
	var time1, time2 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}

	iter, err := ExecuteParallel(NewTask(t1, NewImmediateExecutor()), NewTask(t2, NewImmediateExecutor())).Await(context.Background())
	results := iter.Next()
	assertNil(t, err)
	assertEqual(t, 2, len(results))
	assertNotNil(t, time1)
	assertNotNil(t, time2)
	assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
	assertNil(t, results[0])
	assertNil(t, results[1])
}

func TestExecuteParallel_Error(t *testing.T) {
	var time1, time3 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		return errors.New("test")
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		return errors.New("test")
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		iter, err := ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i]), NewTask(t3, executors[i]), NewTask(t4, executors[i])).Await(context.Background())
		results := iter.Next()
		assertNotNil(t, err)
		assertEqual(t, 4, len(results))
		assertNotNil(t, time1)
		assertNotNil(t, time3)
		assertNil(t, results[0])
		assertNil(t, results[2])
		assertNotNil(t, results[1])
		assertNotNil(t, results[3])
	}
}

func TestCascade_SerialAfterParallel(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return nil
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
				assertNil(t, results[0])
				assertNil(t, results[1])
			} else if i == 1 {
				assertNil(t, results[0])
				assertNil(t, results[1])
			}
			i++
		}
		assertEqual(t, 2, i)
	}
}

func TestCascade_ParallelAfterSerial(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return nil
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
				assertNil(t, results[0])
				assertNil(t, results[1])
			} else if i == 1 {
				assertNil(t, results[0])
				assertNil(t, results[1])
			}
			i++
		}
		assertEqual(t, 2, i)
	}
}

func TestCascade_SerialAfterParallelError(t *testing.T) {
	var time1, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		return errors.New("test")
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return nil
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
			assertNil(t, results[0])
			assertNotNil(t, results[1])
			i++
		}
		assertEqual(t, 1, i)
	}
}

func TestCascade_SerialAfterParallelAfterParallel(t *testing.T) {
	var time1, time2, time3, time4, time5, time6 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return nil
	}
	var t5 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time5 = &now
		return nil
	}
	var t6 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time6 = &now
		return nil
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
				assertNil(t, results[0])
				assertNil(t, results[1])
			} else if i == 1 {
				assertNil(t, results[0])
				assertNil(t, results[1])
			} else if i == 2 {
				assertNil(t, results[0])
				assertNil(t, results[1])
			}
			i++
		}
		assertEqual(t, 3, i)
	}
}

func TestExecution_Async(t *testing.T) {
	var time1, time2, time3, time4 *time.Time
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return nil
	}
	executors := []TaskExecutor{NewPoolExecutor(20, 20), NewAsyncExecutor()}
	for i := range executors {
		wg := sync.WaitGroup{}
		wg.Add(1)
		ExecuteParallel(NewTask(t1, executors[i]), NewTask(t2, executors[i])).
			ExecuteSerial(NewTask(t3, executors[i]), NewTask(t4, executors[i])).
			Async(context.Background(), executors[i], func(iter ExecutionResultIterator, err error) {
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
						assertNil(t, results[0])
						assertNil(t, results[1])
					} else if i == 1 {
						assertNil(t, results[0])
						assertNil(t, results[1])
					}
					i++
				}
				assertEqual(t, 2, i)
				wg.Done()
			})
		wg.Wait()
	}
}

func TestExecution_Switch(t *testing.T) {
	var taskResult string
	var defaultCase = NewTask(func(ctx context.Context) error {
		taskResult = "default"
		return nil
	}, NewAsyncExecutor())

	var case1 = NewTask(func(ctx context.Context) error {
		taskResult = "case1"
		return nil
	}, NewAsyncExecutor())

	var case2 = NewTask(func(ctx context.Context) error {
		taskResult = "case2"
		return nil
	}, NewAsyncExecutor())

	_, _ = Switch(defaultCase.ToExecution(),
		CaseExecution{
			Execution: case1.ToExecution(),
			Case: func() bool {
				return true
			},
		},
		CaseExecution{
			Execution: case2.ToExecution(),
			Case: func() bool {
				return false
			},
		}).Await(context.Background())

	assertEqual(t, taskResult, "case1")

	var defaultCase1 = NewTask(func(ctx context.Context) error {
		taskResult = "defaultCase1"
		return nil
	}, NewAsyncExecutor())

	var case3 = NewTask(func(ctx context.Context) error {
		taskResult = "case3"
		return nil
	}, NewAsyncExecutor())

	_, _ = Switch(defaultCase1.ToExecution(), CaseExecution{
		Execution: case3.ToExecution(),
		Case: func() bool {
			return false
		},
	}).Await(context.Background())

	assertEqual(t, taskResult, "defaultCase1")

	var firstExec = NewTask(func(ctx context.Context) error {
		return nil
	}, NewAsyncExecutor())

	_, _ = ExecuteSerial(firstExec).Switch(defaultCase.ToExecution(),
		CaseExecution{
			Execution: case1.ToExecution(),
			Case: func() bool {
				return true
			},
		},
		CaseExecution{
			Execution: case2.ToExecution(),
			Case: func() bool {
				return false
			},
		}).Await(context.Background())

	assertEqual(t, taskResult, "case1")

	_, _ = ExecuteSerial(firstExec).Switch(defaultCase1.ToExecution(), CaseExecution{
		Execution: case3.ToExecution(),
		Case: func() bool {
			return false
		},
	}).Await(context.Background())

	assertEqual(t, taskResult, "defaultCase1")
}

