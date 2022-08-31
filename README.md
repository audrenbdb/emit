Emit provides a generic emitter that one can subscribe to.

# Usage

Given one wants to emit a list of int.

```go
  ctx := context.Background()
  emitter := emit.NewEmitter[int](ctx)
  emitter.OnEmit(ctx, func(n int) {
    fmt.Printf("Received number: %d\n", n)
  })

  for _, n := range []int{1, 5, 8, 10} {
    emitter.Emit(n)
  }

  // output:
  // Received number: 1
  // Received number: 5
  // Received number: 8
  // Received number: 10
```

# Options

Listened events with OnEmit can be filtered with a boolean match.

Only matching events will be emitted.

I.E: 
```go
  match := func(n int) bool { return n > 5 }
  emitter.OnEmit(ctx, func(n int) {
    fmt.Printf("Received number: %d\n", n)
  })

  for _, n := range []int{1, 5, 8, 10} {
    emitter.Emit(n)
  }
  // output:
  // Received number: 8
  // Received number: 10
```


