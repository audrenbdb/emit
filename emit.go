// Package emit provides a generic event emitter that one can
// subscribe to with OnEmit.
package emit

import (
	"context"
)

type emitter[T any] struct {
	broadcast    chan T
	actions      []action[T]
	addAction    chan action[T]
	removeAction chan action[T]
}

type action[T any] func(event T)

// emitterOption modifies emitter default behavior
type emitterOption[T any] func(em *emitter[T])

// WithBuffer adds a buffer to the emitter
func WithBuffer[T any](bufferSize int) emitterOption[T] {
	return func(em *emitter[T]) {
		em.broadcast = make(chan T, bufferSize)
	}
}

// NewEmitter creates a new event emitter.
// Event emitter can register new actions triggered on emit with OnEmit.
//
// Closing context shuts down the emitter.
func NewEmitter[T any](ctx context.Context, opts ...emitterOption[T]) *emitter[T] {
	em := &emitter[T]{
		broadcast:    make(chan T),
		actions:      make([]action[T], 0),
		addAction:    make(chan action[T]),
		removeAction: make(chan action[T]),
	}
	for _, opt := range opts {
		opt(em)
	}
	go em.run(ctx)
	return em
}

func (em *emitter[T]) run(ctx context.Context) {
	for {
		select {
		case event := <-em.broadcast:
			for _, actionFn := range em.actions {
				actionFn(event)
			}
		case action := <-em.addAction:
			em.actions = append(em.actions, action)
		case action := <-em.removeAction:
			if i := em.getActionIndex(action); i != -1 {
				em.actions[i] = em.actions[len(em.actions)-1]
				em.actions = em.actions[:len(em.actions)-1]
			}
		case <-ctx.Done():
			close(em.broadcast)
			close(em.addAction)
			close(em.removeAction)
			return
		}
	}
}

func (em *emitter[T]) getActionIndex(action action[T]) int {
	for i, a := range em.actions {
		if &a == &action {
			return i
		}
	}
	return -1
}

// Emit broadcasts event to matching listeners
func (em *emitter[T]) Emit(event T) { em.broadcast <- event }

// OnEmit registers a new action triggered given an event is emitted if all conditions match.
//
// Cancelling the context prevent further event to be processed.
func (em *emitter[T]) OnEmit(ctx context.Context, action action[T], conditions ...condition[T]) {
	matchingAction := func(event T) {
		for _, match := range conditions {
			if !match(event) {
				return
			}
		}
		action(event)
	}

	go func() {
		<-ctx.Done()
		em.removeAction <- matchingAction
	}()

	em.addAction <- matchingAction
}

// condition to be true in order to process event
type condition[T any] func(event T) bool
