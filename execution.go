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
	results      [][]*TaskResult
}

func (it *ExecutionResultIterator) Next() []*TaskResult {
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

func (e *Execution) parallelExec(ctx context.Context) ([]*TaskResult, error) {
	var futures []TaskFuture
	var ret []*TaskResult
	for i := range e.tasks {
		futures = append(futures, e.tasks[i].executor.Execute(ctx, e.tasks[i].taskFunc))
	}
	var err error
	for i := range futures {
		result := futures[i].Get()
		if result.Error != nil {
			if err == nil {
				err = result.Error
			} else {
				err = fmt.Errorf("%s:%w", result.Error, err)
			}
		}
		ret = append(ret, result)
	}
	return ret, err
}

func (e *Execution) serialExec(ctx context.Context) ([]*TaskResult, error) {
	var ret []*TaskResult
	for i := range e.tasks {
		result := e.tasks[i].executor.Execute(ctx, e.tasks[i].taskFunc).Get()
		ret = append(ret, result)
		if result.Error != nil {
			return ret, result.Error
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

func (e *Execution) Async(ctx context.Context, executor TaskExecutor, callback func(*ExecutionResultIterator, error)) {
	var execTaskFunc TaskFunc = func(ctx context.Context) *TaskResult {
		result, err := e.Await(ctx)
		if callback != nil {
			callback(result, err)
		}
		return NewTaskResult(result, err)
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
		var execResult []*TaskResult
		switch iter.execType {
		case executionTypeParallel:
			execResult, err = iter.parallelExec(ctx)
		default:
			execResult, err = iter.serialExec(ctx)
		}
		ret.results = append(ret.results, execResult)
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
