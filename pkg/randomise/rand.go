package randomise

import (
	"context"
	"math/rand"

	"github.com/tilinna/clock"
)

type contextKey int

const (
	_ contextKey = iota
	randomCtx
	sourceCtx
)

// FromContext will check if there is an existing
// number generator as part of the context and if it is
// it will then use that.
// Then check if there is a source to use and return
// a new random from that source.
// Otherwise a new randomised source is created and used
func FromContext(ctx context.Context) *rand.Rand {
	if r, ok := ctx.Value(randomCtx).(*rand.Rand); ok {
		return r
	}
	if s, ok := ctx.Value(sourceCtx).(rand.Source); ok {
		return rand.New(s)
	}
	return rand.New(rand.NewSource(clock.FromContext(ctx).Now().Unix()))
}

// WithRandSource adds the Rand.Source as part of the return context
func WithRandSource(ctx context.Context, s rand.Source) context.Context {
	return context.WithValue(ctx, sourceCtx, s)
}

// WithRandom sets rand.Rand as part of the returned context
func WithRandom(ctx context.Context, r *rand.Rand) context.Context {
	return context.WithValue(ctx, randomCtx, r)
}
