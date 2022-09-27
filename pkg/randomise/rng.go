package randomise

import (
	"context"
	"math/rand"

	"github.com/tilinna/clock"
)

type fixed struct {
	value int64
}

var (
	_ rand.Source = (*fixed)(nil)
)

func (f *fixed) Int63() int64 { return f.value }
func (f *fixed) Seed(_ int64) { /* ignored */ }

// NewFixedSource will always return a random number of 4
// there is no alternatives
func NewFixedSource(_ context.Context) rand.Source {
	return &fixed{value: 4}
}

// NewRepeatingSource returns a sequence that will always
// be the same for each calling.
func NewRepeatingSource(_ context.Context) rand.Source {
	return rand.New(rand.NewSource(0))
}

// NewTimeSource uses time to seed the random number generator
// So that each invocation should appear random
func NewTimeSource(ctx context.Context) rand.Source {
	return rand.New(rand.NewSource(clock.FromContext(ctx).Now().Unix()))
}
