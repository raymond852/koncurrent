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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Immediate(), t2.Immediate()},
		{t1.Async(), t2.Async()},
		{t1.Pool(pe), t2.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteSerial(taskExecs[i]...).Await(context.Background())
		results := iter[0]
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Immediate(), t2.Immediate(), t3.Immediate()},
		{t1.Async(), t2.Async(), t3.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteSerial(taskExecs[i]...).Await(context.Background())
		resultErrors := iter[0]
		assertNotNil(t, err)
		assertNil(t, resultErrors[2])
		assertNotNil(t, time1)
		assertNil(t, time3)
		assertNotNil(t, resultErrors[1])
	}
}

func TestExecuteParallel(t *testing.T) {
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async()},
		{t1.Pool(pe), t2.Pool(pe)},
	}
	for i := range taskExecs {
		now := time.Now()
		iter, err := ExecuteParallel(taskExecs[i]...).Await(context.Background())
		elapse := time.Now().Sub(now)
		assertTrue(t, elapse < 1100*time.Millisecond)
		resultErrors := iter[0]
		assertNil(t, err)
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

	iter, err := ExecuteParallel(t1.Immediate(), t2.Immediate()).Await(context.Background())
	results := iter[0]
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteParallel(taskExecs[i]...).Await(context.Background())
		results := iter[0]
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteSerial(taskExecs[i][:2]...).
			ExecuteParallel(taskExecs[i][2:]...).
			Await(context.Background())
		assertNil(t, err)
		for _, results := range iter {
			assertEqual(t, 2, len(results))
			assertNotNil(t, time1)
			assertNotNil(t, time2)
			assertTrue(t, (*time2).Sub(*time1) >= 100*time.Millisecond)
			assertNotNil(t, time3)
			assertTrue(t, (*time3).Sub(*time2) >= 100*time.Millisecond)
			assertNotNil(t, time4)
			assertTrue(t, (*time4).Sub(*time3) < 10*time.Millisecond)
			assertNil(t, results[0])
			assertNil(t, results[1])
		}
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteParallel(taskExecs[i][:2]...).
			ExecuteSerial(taskExecs[i][2:]...).
			Await(context.Background())
		assertNil(t, err)
		i := 0
		for _, results := range iter {
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteSerial(taskExecs[i][:2]...).
			ExecuteParallel(taskExecs[i][2:]...).
			Await(context.Background())
		assertNotNil(t, err)
		assertEqual(t, 2, len(iter[0]))
		assertNil(t, iter[1])
		assertNotNil(t, time1)
		assertNil(t, time3)
		assertNil(t, time4)
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async(), t5.Async(), t6.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe), t5.Pool(pe), t6.Pool(pe)},
	}
	for i := range taskExecs {
		iter, err := ExecuteSerial(taskExecs[i][:2]...).
			ExecuteParallel(taskExecs[i][2:4]...).
			ExecuteParallel(taskExecs[i][4:]...).
			Await(context.Background())
		assertNil(t, err)
		i := 0
		for _, results := range iter {
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
	pe := NewPoolExecutor(10, 10)
	taskExecs := [][]TaskExecution{
		{t1.Async(), t2.Async(), t3.Async(), t4.Async()},
		{t1.Pool(pe), t2.Pool(pe), t3.Pool(pe), t4.Pool(pe)},
	}
	for i := range taskExecs {
		wg := sync.WaitGroup{}
		wg.Add(1)
		ExecuteParallel(taskExecs[i][:2]...).
			ExecuteSerial(taskExecs[i][2:]...).
			Async(context.Background(), func(iter ExecutionResults, err error) {
				assertNil(t, err)
				i := 0
				for _, results := range iter {
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

func TestExecuteParallel_WithCancel(t *testing.T) {
	outer := 1
	var task1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		outer = 2
		return nil
	}

	var task2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	var task3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	begin := time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	_, _ = ExecuteParallel(task1.Pool(pe), task2.Async(), task3.Async()).Await(ctx)
	elapsed := time.Now().Sub(begin)
	assertEqual(t, outer, 1)
	assertTrue(t, elapsed < 1000*time.Millisecond)
}

func TestExecuteSerial_WithCancel(t *testing.T) {
	outer := 1
	var task1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		outer = 2
		return nil
	}

	var task2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	var task3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	begin := time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	_, _ = ExecuteSerial(task1.Pool(pe), task2.Async(), task3.Async()).Await(ctx)
	elapsed := time.Now().Sub(begin)
	assertEqual(t, outer, 1)
	assertTrue(t, elapsed < 1000*time.Millisecond)
}

func TestExecution_Switch(t *testing.T) {
	var taskResult string
	var defaultCase TaskFunc = func(ctx context.Context) error {
		taskResult = "default"
		return nil
	}

	var case1 TaskFunc = func(ctx context.Context) error {
		taskResult = "case1"
		return nil
	}

	var case2 TaskFunc = func(ctx context.Context) error {
		taskResult = "case2"
		return nil
	}

	_, _ = Switch(defaultCase.Async().Execution(),
		CaseExecution{
			Execution: case1.Async().Execution(),
			Case: func() bool {
				return true
			},
		},
		CaseExecution{
			Execution: case2.Async().Execution(),
			Case: func() bool {
				return false
			},
		}).Await(context.Background())

	assertEqual(t, taskResult, "case1")

	var defaultCase1 TaskFunc = func(ctx context.Context) error {
		taskResult = "defaultCase1"
		return nil
	}

	var case3 TaskFunc = func(ctx context.Context) error {
		taskResult = "case3"
		return nil
	}

	_, _ = Switch(defaultCase1.Async().Execution(), CaseExecution{
		Execution: case3.Async().Execution(),
		Case: func() bool {
			return false
		},
	}).Await(context.Background())

	assertEqual(t, taskResult, "defaultCase1")

	var firstExec TaskFunc = func(ctx context.Context) error {
		return nil
	}

	_, _ = firstExec.Async().Execution().
		Switch(defaultCase.Async().Execution(),
			CaseExecution{
				Execution: case1.Async().Execution(),
				Case: func() bool {
					return true
				},
			},
			CaseExecution{
				Execution: case2.Async().Execution(),
				Case: func() bool {
					return false
				},
			}).Await(context.Background())

	assertEqual(t, taskResult, "case1")

	_, _ = firstExec.Async().Execution().
		Switch(defaultCase1.Async().Execution(),
			CaseExecution{
				Execution: case3.Async().Execution(),
				Case: func() bool {
					return false
				},
			}).Await(context.Background())

	assertEqual(t, taskResult, "defaultCase1")
}
