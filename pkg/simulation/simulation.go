package simulation

import "context"

// Simulation represents an environment that
// experiments are taking place and advances
// the state by one each call
type Simulation[Result any] interface {
	Next(ctx context.Context) Result
}

// Simulation function lets you run an existing function
// as part of the simulation
type SimulationFunc[T any] func(ctx context.Context) T

var (
	_ Simulation[any] = (SimulationFunc[any])(nil)
)

func (fn SimulationFunc[T]) Next(ctx context.Context) T {
	return fn(ctx)
}
