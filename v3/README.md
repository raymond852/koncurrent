## Introduction
A Go lib for easier concurrency control. Inspired by ReactiveX and javascript Promise.

The v3 library dramatically improve performance than v2 and v1 library by reducing the heap memory allocation.

### Usage
#### Simple execution example
```go
    var time1, time2 time.Time
    var t1 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        time1 = time.Now()
        return nil
    }
    var t2 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        time2 = time.Now()
        return nil
    }
    errIter, err := koncurrent.ExecuteParallel(t1.Async(), t2.Immediate()).Await(context.Background())
    fmt.Println(errIter)
    fmt.Println(err)
```
#### Cascaded execution example
```go
    var time1, time2, time3, time4 time.Time
    var t1 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        time1 = time.Now()
        return nil
    }
    var t2 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        time2 = time.Now()
        return nil
    }
    var t3 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        tim3 = time.Now()
        return nil
    }
    var t4 koncurrent.TaskFunc = func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        time4 = time.Now()
        return errors.New("task 4 error occur")
    }
    pe := koncurrent.NewPoolExecutor(20, 20)
    for i := range executors {
        errIter, err := koncurrent.ExecuteSerial(t1.Pool(pe), t2.Async()).
            ExecuteParallel(t3.Pool(pe), t4.Pool(pe)).
            Await(context.Background())
        fmt.Println(errIter)
        fmt.Println(err)
    }
```
#### Check more example in execution_test.go