# Event Emitter

## Description

A generic event emitter.

## Usage

```go
em := emit.New[int](ctx)
numbers := make([]int, 3)
sub := em.Sub(ctx)

em.Emit(1, 2, 3)

for i := 0; i < 3; i++ {
    numbers[i] = <-sub
}

fmt.Println(numbers)
// Output: [1, 2, 3]
```