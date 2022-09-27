package simulation

import (
	"context"

	"github.com/MovieStoreGuy/datagen/pkg/dice"
	"github.com/MovieStoreGuy/datagen/pkg/generic"
)

// NewDiceSimulation uses a collection of dice which and then each rolled
// on each Simulation function
func NewDiceSimulation(dies []dice.Die) Simulation[[]int] {
	return SimulationFunc[[]int](func(ctx context.Context) []int {
		return generic.ParallelEachFuncWithContext(ctx, generic.StreamMap(func(d dice.Die) func(context.Context) int {
			return d.Roll
		}, dies...)...)
	})
}
