## Introduction
A Go lib for easier concurrency control. Inspired by ReactiveX and javascript Promise.

The v3 library dramatically improve performance than v2 and v1 library by reducing the heap memory allocation.


### Benchmark result
Benchmark test on an AMD 2700X Ubuntu 20.04.3 LTS machine Go 1.16.8
#### v1
```
BenchmarkExecuteSerial_Immediate-16      	 1259301	       926.5 ns/op	     304 B/op	      14 allocs/op
BenchmarkExecuteSerial_Async-16          	  317821	      4144 ns/op	     616 B/op	      20 allocs/op
BenchmarkExecuteSerial_Pool-16           	  312390	      3894 ns/op	     616 B/op	      20 allocs/op
BenchmarkExecuteParallel_Immediate-16    	  897908	      1252 ns/op	     416 B/op	      17 allocs/op
BenchmarkExecuteParallel_Async-16        	  252895	      5007 ns/op	     728 B/op	      23 allocs/op
BenchmarkExecuteParallel_Pool-16         	  254157	      4411 ns/op	     728 B/op	      23 allocs/op
```
#### original v2
```
BenchmarkExecuteSerial_Immediate-16      	 1259301	       926.5 ns/op	     304 B/op	      14 allocs/op
BenchmarkExecuteSerial_Async-16          	  317821	      4144 ns/op	     616 B/op	      20 allocs/op
BenchmarkExecuteSerial_Pool-16           	  312390	      3894 ns/op	     616 B/op	      20 allocs/op
BenchmarkExecuteParallel_Immediate-16    	  897908	      1252 ns/op	     416 B/op	      17 allocs/op
BenchmarkExecuteParallel_Async-16        	  252895	      5007 ns/op	     728 B/op	      23 allocs/op
BenchmarkExecuteParallel_Pool-16         	  254157	      4411 ns/op	     728 B/op	      23 allocs/op
```
#### optimized v2 branch
```
BenchmarkExecuteSerial_Immediate-16              2785822               421.7 ns/op           120 B/op          5 allocs/op
BenchmarkExecuteSerial_Async-16                   380536              3064 ns/op             408 B/op          8 allocs/op
BenchmarkExecuteSerial_Pool-16                    477312              2824 ns/op             408 B/op          8 allocs/op
BenchmarkExecuteParallel_Immediate-16            2117505               543.6 ns/op           168 B/op          6 allocs/op
BenchmarkExecuteParallel_Async-16                 348438              3392 ns/op             456 B/op          9 allocs/op
BenchmarkExecuteParallel_Pool-16                  414778              2790 ns/op             456 B/op          9 allocs/op
```
#### v3
```
BenchmarkExecuteSerial_Immediate-16              4078363               286.0 ns/op           168 B/op          3 allocs/op
BenchmarkExecuteSerial_Async-16                   432781              2540 ns/op             168 B/op          3 allocs/op
BenchmarkExecuteSerial_Pool-16                    452028              2473 ns/op             168 B/op          3 allocs/op
BenchmarkExecuteParallel_Immediate-16            1586167               762.7 ns/op           248 B/op          4 allocs/op
BenchmarkExecuteParallel_Async-16                 384945              3181 ns/op             248 B/op          4 allocs/op
BenchmarkExecuteParallel_Pool-16                  425948              2701 ns/op             248 B/op          4 allocs/op
```

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