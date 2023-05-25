package emit_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/audrenbdb/emit"
)

func TestEmitter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	t.Cleanup(cancel)

	t.Run("EmitOneSub", func(t *testing.T) {
		em := emit.New[int](ctx)
		sub := em.Sub(ctx)

		em.Emit(1, 2, 3)

		numbers := make([]int, 3)
		for i := range numbers {
			select {
			case <-ctx.Done():
				t.Errorf("timeout")
			case n := <-sub:
				numbers[i] = n
			}
		}

		if diff := cmp.Diff(numbers, []int{1, 2, 3}); diff != "" {
			t.Errorf("numbers mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("EmitTwoSub", func(t *testing.T) {
		em := emit.New[int](ctx)
		firstSub := em.Sub(ctx)
		secondSub := em.Sub(ctx)

		em.Emit(1)

		select {
		case <-ctx.Done():
			t.Errorf("timeout")
		case n := <-firstSub:
			if n != 1 {
				t.Errorf("firstSub mismatch (-want +got):\n%s", cmp.Diff(n, 1))
			}
		}

		select {
		case <-ctx.Done():
			t.Errorf("timeout")
		case n := <-secondSub:
			if n != 1 {
				t.Errorf("secondSub mismatch (-want +got):\n%s", cmp.Diff(n, 1))
			}
		}
	})

	t.Run("EmitOneSubWithMatch", func(t *testing.T) {
		em := emit.New[int](ctx)
		sub := em.Sub(ctx, func(i int) bool {
			return i > 5
		})

		em.Emit(1, 2, 3, 4, 5, 6)

		select {
		case <-ctx.Done():
			t.Errorf("timeout")
		case n := <-sub:
			if n != 6 {
				t.Errorf("only 6 should be received")
			}
		}
	})
}
