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
	tasksList         [][]TaskExecution
	executionTypeList []int
}

type CaseExecution struct {
	Execution Execution
	Case      func() bool
}

func (e Execution) nextExecution(tasks []TaskExecution, executionType int) Execution {
	return Execution{
		tasksList:         append(e.tasksList, tasks),
		executionTypeList: append(e.executionTypeList, executionType),
	}
}

func (e Execution) ExecuteParallel(tasks ...TaskExecution) Execution {
	return e.nextExecution(tasks, executionTypeParallel)
}

func (e Execution) ExecuteSerial(tasks ...TaskExecution) Execution {
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

func (e Execution) Async(ctx context.Context, callback func(ExecutionResults, error)) {
	go func() {
		result, err := e.Await(ctx)
		if callback != nil {
			callback(result, err)
		}
	}()
}

func (e Execution) Await(ctx context.Context) (ExecutionResults, error) {
	var ret ExecutionResults = make([][]error, len(e.tasksList))
	var err error
	for i := range e.tasksList {
		currTaskList := e.tasksList[i]
		execErr := make([]error, len(currTaskList))
		ret[i] = execErr
		switch e.executionTypeList[i] {
		case executionTypeParallel:
			futures := make([]TaskFuture, len(currTaskList))
			for j, task := range currTaskList {
				if task.executionType == taskExecutionTypeImmediate {
					futures[j] = make(chan error, 1)
					futures[j] <- task.taskFunc(ctx)
				} else if task.executionType == taskExecutionTypeNewGoRoutine {
					chn := make(chan error)
					taskFunc :=  task.taskFunc
					go func() {
						chn <- taskFunc(ctx)
					}()
					futures[j] = chn
				} else {
					futures[j] = task.executor.Execute(ctx, task.taskFunc)
				}
			}
			for j := range futures {
				resultErr := <- futures[j]
				close(futures[j])
				if err == nil {
					err = resultErr
				} else {
					err = fmt.Errorf("%s:%w", resultErr, err)
				}
				execErr[j] = resultErr
			}
		default:
			chn := make(chan error)
			for j, task := range currTaskList {
				var resultErr error
				if task.executionType == taskExecutionTypeImmediate {
					resultErr = task.taskFunc(ctx)
				} else if task.executionType == taskExecutionTypeNewGoRoutine {
					taskFunc :=  task.taskFunc
					go func() {
						chn <- taskFunc(ctx)
					}()
					resultErr = <-chn
				} else {
					resultErr = <-task.executor.Execute(ctx, task.taskFunc)
				}
				execErr[j] = resultErr
				if resultErr != nil {
					close(chn)
					return ret, resultErr
				}
			}
			close(chn)
		}
		if err != nil {
			break
		}
	}
	return ret, err
}

func ExecuteParallel(tasks ...TaskExecution) Execution {
	return Execution{
		tasksList:         [][]TaskExecution{tasks},
		executionTypeList: []int{executionTypeParallel},
	}
}

func ExecuteSerial(tasks ...TaskExecution) Execution {
	return Execution{
		tasksList:         [][]TaskExecution{tasks},
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
