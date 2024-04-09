package await

import (
	"context"
	"sync"
)

type contextKey struct{}

var awaitKey = &contextKey{}

// Awaiter basically an interface that describes the sync.WaitGroup.
type Awaiter interface {
	Add(delta int)
	Done()
	Wait()
}

// nopAwaiter is a no-op implementation of Awaiter.
type nopAwaiter struct{}

func (nopAwaiter) Add(_ int) {}
func (nopAwaiter) Done()     {}
func (nopAwaiter) Wait()     {}

// Context returns a new context with an Awaiter.
func Context(ctx context.Context) context.Context {
	var wg Awaiter = &sync.WaitGroup{}
	return context.WithValue(ctx, awaitKey, wg)
}

// FromContext returns the Awaiter from the context if it exists. Otherwise, it returns a no-op Awaiter.
func FromContext(ctx context.Context) Awaiter {
	wg, ok := ctx.Value(awaitKey).(Awaiter)
	if !ok {
		return &nopAwaiter{}
	}
	return wg
}
