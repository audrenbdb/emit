package emit

import (
	"context"
)

// Emitter emits an event.
type Emitter[T any] struct {
	sub   chan *subscriber[T]
	unsub chan *subscriber[T]
	emit  chan T
}

type subscriber[T any] struct {
	ch       chan T
	matchers []func(T) bool
}

// New constructs a new [Emitter].
func New[T any](ctx context.Context) *Emitter[T] {
	emitter := &Emitter[T]{
		sub:   make(chan *subscriber[T]),
		unsub: make(chan *subscriber[T]),
		emit:  make(chan T, 1024),
	}

	go func() {
		subs := map[*subscriber[T]]bool{}

		for {
			select {
			case <-ctx.Done():
				for sub := range subs {
					close(sub.ch)
					delete(subs, sub)
				}

				return
			case sub := <-emitter.sub:
				subs[sub] = true
			case sub := <-emitter.unsub:
				if _, ok := subs[sub]; ok {
					delete(subs, sub)
					close(sub.ch)
				}
			case event := <-emitter.emit:
				for sub := range subs {
					if !emitter.checkMatchers(event, sub.matchers...) {
						continue
					}

					select {
					case sub.ch <- event:
					default:
						// channel has been closed on consumer side / or deadlock => close.
						close(sub.ch)
						delete(subs, sub)
					}
				}
			}
		}
	}()

	return emitter
}

// checkMatchers checks that event matches all matchers.
func (e *Emitter[T]) checkMatchers(event T, matchers ...func(T) bool) bool {
	for _, match := range matchers {
		if !match(event) {
			return false
		}
	}

	return true
}

// Emit emits given events.
func (e *Emitter[T]) Emit(events ...T) {
	for _, event := range events {
		e.emit <- event
	}
}

// Sub subscribes to all events. Closes subscription on context done. Only matching events are sent to the chan.
func (e *Emitter[T]) Sub(ctx context.Context, matchers ...func(T) bool) chan T {
	sub := &subscriber[T]{ch: make(chan T, 1024), matchers: matchers}

	go func() {
		<-ctx.Done()
		e.unsub <- sub
	}()

	e.sub <- sub

	return sub.ch
}
