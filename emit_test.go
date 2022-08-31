package emit_test

import (
	"context"
	"github.com/audrenbdb/emit"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	t.Run("Only positive numbers should be emitted", func(t *testing.T) {
		ctx := context.Background()

		emitter := emit.NewEmitter[int](ctx)
		match := func(n int) bool { return n > 0 }

		var numbers []int

		emitter.OnEmit(ctx, func(n int) {
			numbers = append(numbers, n)
		}, match)

		for _, n := range []int{-1, 0, -3, 1, -100000, 2, 3} {
			emitter.Emit(n)
		}

		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				if cmp.Equal(numbers, []int{1, 2, 3}) {
					return
				}
			case <-time.After(100 * time.Millisecond):
				t.Errorf("expected numbers 1, 2, 3, got: %v", numbers)
				return
			}
		}
	})
}
