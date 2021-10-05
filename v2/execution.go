package koncurrent

import (
	"context"
	"fmt"
)

const (
	executionTypeParallel = iota
	executionTypeSerial
)

type ExecutionResults [][]error

type Execution struct {
	tasksList         [][]Task
	executionTypeList []int
}

type CaseExecution struct {
	Execution Execution
	Case      func() bool
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

func (e Execution) Async(ctx context.Context, executor TaskExecutor, callback func(ExecutionResults, error)) {
	var execTaskFunc TaskFunc = func(ctx context.Context) error {
		result, err := e.Await(ctx)
		if callback != nil {
			callback(result, err)
		}
		return err
	}
	_ = executor.Execute(ctx, execTaskFunc)
}

func (e Execution) Await(ctx context.Context) (ExecutionResults, error) {
	ret := make([][]error, len(e.tasksList))

	var err error
	for i := range e.tasksList {
		execErr := make([]error, len(e.tasksList[i]))
		ret[i] = execErr
		switch e.executionTypeList[i] {
		case executionTypeParallel:
			futures := make([]TaskFuture, len(e.tasksList[i]))
			for j, task := range e.tasksList[i] {
				futures[j] = task.executor.Execute(ctx, task.taskFunc)
			}
			for j := range futures {
				resultErr := futures[j].Get()
				execErr[j] = resultErr
				if err == nil {
					err = resultErr
				} else {
					err = fmt.Errorf("%s:%w", resultErr, err)
				}
			}
		default:
			for j, task := range e.tasksList[i] {
				resultErr := task.executor.Execute(ctx, task.taskFunc).Get()
				execErr[j] = resultErr
				if resultErr != nil {
					return ret, resultErr
				}
			}
		}
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
