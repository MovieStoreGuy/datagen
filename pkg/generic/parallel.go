package generic

import (
	"context"
	"runtime"
	"sync"
)

// ParallelEachFuncWithContext is a lock free parallel function that will execute each task
// and store the result. Once all the tasks are complate, the function will then return.
func ParallelEachFuncWithContext[T any](ctx context.Context, tasks ...func(context.Context) T) []T {
	var (
		wg     sync.WaitGroup
		values = make([]T, len(tasks))
		cores  = runtime.NumCPU()
	)
	wg.Add(cores)
	for i := 0; i < cores; i++ {
		go func(ctx context.Context, offset int) {
			defer wg.Done()

			for i := offset; i < len(tasks); i += cores {
				values[i] = tasks[i](ctx)
			}
		}(ctx, i)
	}
	wg.Wait()
	return values
}
