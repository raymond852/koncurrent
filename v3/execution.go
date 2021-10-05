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
			resultsChn := make(chan TaskResult, len(currTaskList))
			for j, task := range currTaskList {
				taskFunc := task.taskFunc
				if task.executionType == taskExecutionTypeImmediate {
					resultsChn <- TaskResult{
						err: taskFunc(ctx),
						id:  j,
					}
				} else if task.executionType == taskExecutionTypeNewGoRoutine {
					index := j
					go func() {
						resultsChn <- TaskResult{
							err: taskFunc(ctx),
							id:  index,
						}
					}()
				} else {
					task.executor.Execute(ctx, task.taskFunc, j, resultsChn)
				}
			}
			for range currTaskList {
				select {
				case taskResult := <-resultsChn:
					execErr[taskResult.id] = taskResult.err
				case <-ctx.Done():
					return ret, err
				}
			}
			close(resultsChn)
			for j := range execErr {
				if execErr[j] != nil {
					if err == nil {
						err = execErr[j]
					} else {
						err = fmt.Errorf("%s:%w", execErr[j], err)
					}
				}
			}
		default:
			resultsChn := make(chan TaskResult)
			for j, task := range currTaskList {
				taskFunc := task.taskFunc
				if task.executionType == taskExecutionTypeImmediate {
					execErr[j] = taskFunc(ctx)
				} else if task.executionType == taskExecutionTypeNewGoRoutine {
					go func() {
						resultsChn <- TaskResult{
							err: taskFunc(ctx),
						}
					}()
					select {
					case taskResult := <-resultsChn:
						execErr[j] = taskResult.err
					case <-ctx.Done():
						return ret, err
					}
				} else {
					task.executor.Execute(ctx, task.taskFunc, j, resultsChn)
					select {
					case taskResult := <-resultsChn:
						execErr[j] = taskResult.err
					case <-ctx.Done():
						return ret, err
					}
				}
				if execErr[j] != nil {
					close(resultsChn)
					return ret, execErr[j]
				}
			}
			close(resultsChn)
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
