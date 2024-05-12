package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := &sync.WaitGroup{}
	wg.Add(n)

	muTaskIndex := &sync.Mutex{}
	currentTask := 0

	muErrCount := &sync.Mutex{}
	curErrCount := 0

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for {
				muTaskIndex.Lock()
				taskIndex := currentTask
				currentTask++
				muTaskIndex.Unlock()

				if taskIndex >= len(tasks) {
					break
				}

				if err := tasks[taskIndex](); err != nil {
					muErrCount.Lock()
					curErrCount++

					if curErrCount >= m {
						muErrCount.Unlock()
						break
					}
					muErrCount.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if curErrCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
