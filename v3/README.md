## Introduction
A Go lib for easier concurrency control. Inspired by ReactiveX and javascript Promise

### Executor Types
* Immediate executor. Immediate executor will execute the task on current go routine
* Async executor. Async executor will execute the task on new go routine
* Pool executor. Pool executor will spawn a fixed size go routine pool, all the tasks will be executed by the go routine in the pool

### Usage
#### Simple execution example
```go
    var time1, time2 *time.Time
    var t1 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        now := time.Now()
        time1 = &now
        return nil
    }
	var t2 koncurrent.TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
    // although the below tries to execute in parallel, but since immediate executor is being used, so t2 will execute after t1 finished
    koncurrent.ExecuteParallel(koncurrent.NewTask(t1, koncurrent.NewImmediateExecutor()), koncurrent.NewTask(t2, koncurrent.NewImmediateExecutor()))
```
#### Cascaded execution example
```go
	var time1, time2, time3, time4 *time.Time
	var t1 koncurrent.TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time1 = &now
		return nil
	}
	var t2 koncurrent.TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time2 = &now
		return nil
	}
	var t3 koncurrent.TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time3 = &now
		return nil
	}
	var t4 koncurrent.TaskFunc = func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		now := time.Now()
		time4 = &now
		return errors.New("task 4 error occur")
	}
	executors := []koncurrent.TaskExecutor{koncurrent.NewPoolExecutor(20, 20), koncurrent.NewAsyncExecutor()}
	for i := range executors {
		errIter, err := koncurrent.ExecuteSerial(koncurrent.NewTask(t1, executors[i]), koncurrent.NewTask(t2, executors[i])).
			ExecuteParallel(koncurrent.NewTask(t3, executors[i]), koncurrent.NewTask(t4, executors[i])).
			Await(context.Background())
		fmt.Println(errIter)
		fmt.Println(err)
	}
```