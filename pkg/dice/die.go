package dice

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/MovieStoreGuy/datagen/pkg/randomise"
)

// A die represents a singular dice that can be rolled to give a value
type Die interface {
	// Roll will suspend the dice roll for a certain amount of time
	// then produce a random die value, the value returned will be in the
	// set of [1, Sides]
	Roll(ctx context.Context) int

	Sides() int
}

type die struct {
	sides int
}

var (
	_ Die = (*die)(nil)
)

// NewDie returns a die with the number of sides
func NewDie(sides int) Die {
	return &die{sides: abs(sides)}
}

func (d *die) Sides() int { return d.sides }

func (d *die) Roll(ctx context.Context) (value int) {
	// Since IntN can return a value from [0,d.sides),
	// it needs to be shifted by one to be [1,d.sides]
	return randomise.FromContext(ctx).Intn(d.sides) + 1
}

func abs[N constraints.Integer](n N) N {
	if n < N(0) {
		return -n
	}
	return n
}
