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

func (it *ExecutionResultIterator) Next() []error {
	if it.currentIndex >= len(it.results) {
		return nil
	}
	ret := it.results[it.currentIndex]
	it.currentIndex = it.currentIndex + 1
	return ret
}

type Execution struct {
	root     *Execution
	next     *Execution
	tasks    []*Task
	execType int
}

type CaseExecution struct {
	Execution *Execution
	Case func() bool
}

func (e *Execution) parallelExec(ctx context.Context) ([]error, error) {
	var futures []TaskFuture
	var ret []error
	for i := range e.tasks {
		futures = append(futures, e.tasks[i].executor.Execute(ctx, e.tasks[i].taskFunc))
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

func (e *Execution) serialExec(ctx context.Context) ([]error, error) {
	var ret []error
	for i := range e.tasks {
		resultErr := e.tasks[i].executor.Execute(ctx, e.tasks[i].taskFunc).Get()
		ret = append(ret, resultErr)
		if resultErr != nil {
			return ret, resultErr
		}
	}
	return ret, nil
}

func (e *Execution) nextExecution(tasks []*Task, executionType int) *Execution {
	ret := &Execution{
		root:     nil,
		next:     nil,
		tasks:    tasks,
		execType: executionType,
	}
	if e.root == nil {
		ret.root = e
	} else {
		ret.root = e.root
	}
	e.next = ret
	return ret
}

func (e *Execution) ExecuteParallel(tasks ...*Task) *Execution {
	return e.nextExecution(tasks, executionTypeParallel)
}

func (e *Execution) ExecuteSerial(tasks ...*Task) *Execution {
	return e.nextExecution(tasks, executionTypeSerial)
}

func (e *Execution) Switch(defaultExec *Execution, cases ...CaseExecution) *Execution {
	for i := range cases {
		if cases[i].Case() {
			return cases[i].Execution
		}
	}
	return defaultExec
}

func (e *Execution) Async(ctx context.Context, executor TaskExecutor, callback func(*ExecutionResultIterator, error)) {
	var execTaskFunc TaskFunc = func(ctx context.Context) error {
		result, err := e.Await(ctx)
		if callback != nil {
			callback(result, err)
		}
		return err
	}
	_ = executor.Execute(ctx, execTaskFunc)
}

func (e *Execution) Await(ctx context.Context) (*ExecutionResultIterator, error) {
	iter := e.root
	if e.root == nil {
		iter = e
	}
	ret := &ExecutionResultIterator{}
	var err error
	for {
		var execErr []error
		switch iter.execType {
		case executionTypeParallel:
			execErr, err = iter.parallelExec(ctx)
		default:
			execErr, err = iter.serialExec(ctx)
		}
		ret.results = append(ret.results, execErr)
		if err != nil {
			break
		}
		iter = iter.next
		if iter == nil {
			break
		}
	}
	return ret, err
}

func ExecuteParallel(tasks ...*Task) *Execution {
	return &Execution{
		root:     nil,
		next:     nil,
		tasks:    tasks,
		execType: executionTypeParallel,
	}
}

func ExecuteSerial(tasks ...*Task) *Execution {
	return &Execution{
		root:     nil,
		next:     nil,
		tasks:    tasks,
		execType: executionTypeSerial,
	}
}

func Switch(defaultExec *Execution, cases ...CaseExecution) *Execution {
	for i := range cases {
		if cases[i].Case() {
			return cases[i].Execution
		}
	}
	return defaultExec
}