package koncurrent

import (
	"context"
	"fmt"
)

const (
	executionTypeParallel = iota
	executionTypeSerial
)

type ExecutionResultIterator struct {
	currentIndex int
	results      [][]error
}

func (it ExecutionResultIterator) Next() []error {
	if it.currentIndex >= len(it.results) {
		return nil
	}
	ret := it.results[it.currentIndex]
	it.currentIndex = it.currentIndex + 1
	return ret
}

type Execution struct {
	tasksList         [][]Task
	executionTypeList []int
}

type CaseExecution struct {
	Execution Execution
	Case      func() bool
}

func parallelExec(ctx context.Context, tasks []Task) ([]error, error) {
	var futures []TaskFuture
	var ret []error
	for i := range tasks {
		futures = append(futures, tasks[i].executor.Execute(ctx, tasks[i].taskFunc))
	}
	var err error
	for i := range futures {
		resultErr := futures[i].Get()
		if err == nil {
			err = resultErr
		} else {
			err = fmt.Errorf("%s:%w", resultErr, err)
		}
		ret = append(ret, resultErr)
	}
	return ret, err
}

func serialExec(ctx context.Context, tasks []Task) ([]error, error) {
	var ret []error
	for i := range tasks {
		resultErr := tasks[i].executor.Execute(ctx, tasks[i].taskFunc).Get()
		ret = append(ret, resultErr)
		if resultErr != nil {
			return ret, resultErr
		}
	}
	return ret, nil
}

func (e Execution) nextExecution(tasks []Task, executionType int) Execution {
	return Execution{
		tasksList: append(e.tasksList, tasks),
		executionTypeList: append(e.executionTypeList, executionType),
	}
}

func (e Execution) ExecuteParallel(tasks ...Task) Execution {
	return e.nextExecution(tasks, executionTypeParallel)
}

func (e Execution) ExecuteSerial(tasks ...Task) Execution {
	return e.nextExecution(tasks, executionTypeSerial)
}

func (e Execution) Switch(defaultExec Execution, cases ...CaseExecution) Execution {
	for i := range cases {
		if cases[i].Case() {
			return cases[i].Execution
		}
	}
	return defaultExec
}

func (e Execution) Async(ctx context.Context, executor TaskExecutor, callback func(ExecutionResultIterator, error)) {
	var execTaskFunc TaskFunc = func(ctx context.Context) error {
		result, err := e.Await(ctx)
		if callback != nil {
			callback(result, err)
		}
		return err
	}
	_ = executor.Execute(ctx, execTaskFunc)
}

func (e Execution) Await(ctx context.Context) (ExecutionResultIterator, error) {
	ret := ExecutionResultIterator{}
	var err error
	for i := range e.tasksList {
		var execErr []error
		switch e.executionTypeList[i] {
		case executionTypeParallel:
			execErr, err = parallelExec(ctx, e.tasksList[i])
		default:
			execErr, err = serialExec(ctx, e.tasksList[i])
		}
		ret.results = append(ret.results, execErr)
		if err != nil {
			break
		}
	}
	return ret, err
}

func ExecuteParallel(tasks ...Task) Execution {
	return Execution{
		tasksList: [][]Task{tasks},
		executionTypeList: []int{executionTypeParallel},
	}
}

func ExecuteSerial(tasks ...Task) Execution {
	return Execution{
		tasksList: [][]Task{tasks},
		executionTypeList: []int{executionTypeSerial},
	}
}

func Switch(defaultExec Execution, cases ...CaseExecution) Execution {
	for i := range cases {
		if cases[i].Case() {
			return cases[i].Execution
		}
	}
	return defaultExec
}
