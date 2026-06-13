<div align="center">
    <img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/vowlink.svg)](https://pkg.go.dev/github.com/shengyanli1982/vowlink)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shengyanli1982/vowlink)

## Introduction

`VowLink` is a lightweight Promise library for Go, inspired by the JavaScript Promise API. It replaces callback hell with a fluent, chainable way to handle asynchronous operations.

**Highlights:**

- Chainable API: `Then()` / `Catch()` / `Finally()`
- Full combinators: `All()` / `Race()` / `Any()` / `AllSettled()`
- Zero external dependencies
- Thread-safe via `sync.RWMutex`
- Immutable state: once settled (Fulfilled / Rejected), it never changes

## Installation

```bash
go get github.com/shengyanli1982/vowlink
```

## Quick Start

```go
package main

import (
    "fmt"
    vl "github.com/shengyanli1982/vowlink"
)

func main() {
    result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
        resolve("hello world", nil)
    }).Then(func(value interface{}) (interface{}, error) {
        return value.(string) + " vowlink", nil
    }, nil).Then(func(value interface{}) (interface{}, error) {
        return value.(string) + " !!", nil
    }, nil)

    fmt.Println(result.GetValue()) // hello world vowlink !!
}
```

## Core Concepts

### State Machine

Every Promise transitions through the following states:

```
         resolve(value, nil)
         ┌───────────────────────┐
         │                       ▼
      Pending ─────────────► Fulfilled
         │
         │ reject(nil, err)
         └───────────────────────┐
                                 ▼
                             Rejected
```

- **Pending** — Initial state, waiting to be settled
- **Fulfilled** — Operation succeeded, carrying a return value
- **Rejected** — Operation failed, carrying an error reason

Once a Promise transitions out of Pending, its state is permanent.

### Rules

1. Callbacks in `Then`, `Catch`, and `Finally` all return `(interface{}, error)`. Errors propagate down the chain until caught by `Catch`.
2. `resolve` and `reject` both accept `(value, error)` for maximum flexibility.
3. `GetValue()` and `GetReason()` are **terminal methods** — they return plain values, not Promises.
4. **Do not** create goroutines inside `Then()` / `Catch()` / `Finally()`. For async work, wrap the entire Promise chain in a goroutine.

## API Reference

### Promise Methods

| Method                          | Description                                                         |
| ------------------------------- | ------------------------------------------------------------------- |
| `NewPromise(handler)`           | Creates a Promise; `handler` receives `resolve` and `reject`        |
| `Then(onFulfilled, onRejected)` | Registers success/failure callbacks; chainable; either may be `nil` |
| `Catch(onRejected)`             | Shorthand for `Then(nil, onRejected)`; catches upstream errors      |
| `Finally(cleanup)`              | Runs regardless of outcome; returning an error breaks the chain     |
| `GetValue() interface{}`        | Returns the resolved value (terminal method)                        |
| `GetReason() error`             | Returns the rejection reason (terminal method)                      |

### Combinators

| Method                    | Behavior                                                                            |
| ------------------------- | ----------------------------------------------------------------------------------- |
| `All(p1, p2, ...)`        | Resolves with `[]interface{}` when **all** succeed; rejects on **first** failure    |
| `Race(p1, p2, ...)`       | Settles with the **first** settled Promise (resolve or reject)                      |
| `Any(p1, p2, ...)`        | Resolves with the **first** success; rejects with `*AggregateError` if **all** fail |
| `AllSettled(p1, p2, ...)` | Waits for **all** to settle; collects all values and errors                         |

## Examples

### Chaining & Error Handling

```go
// Success chain
result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    resolve("hello world", nil)
}).Then(func(value interface{}) (interface{}, error) {
    return value.(string) + " vowlink !!", nil
}, nil)

fmt.Println("Resolve:", result.GetValue())
// Resolve: hello world vowlink !!

// Error handling
result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    reject(nil, fmt.Errorf("error"))
}).Then(func(value interface{}) (interface{}, error) {
    return value, nil
}, nil).Catch(func(err error) (interface{}, error) {
    return nil, fmt.Errorf("caught: %s", err.Error())
})

fmt.Println("Rejected:", result.GetReason().Error())
// Rejected: caught: error
```

### Error Recovery

When `Catch` returns a value with a `nil` error, the chain recovers and subsequent `Then` handlers continue executing.

```go
result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    reject(nil, fmt.Errorf("network timeout"))
}).Catch(func(err error) (interface{}, error) {
    return "fallback data", nil // recover
}).Then(func(value interface{}) (interface{}, error) {
    return fmt.Sprintf("got: %v", value), nil
}, nil)

fmt.Println(result.GetValue())
// got: fallback data
```

### Concurrent Collection

```go
p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    resolve("A", nil)
})
p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    resolve("B", nil)
})
p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
    resolve("C", nil)
})

// All — succeeds only when every Promise resolves
result := vl.All(p1, p2, p3)
for i, v := range result.GetValue().([]interface{}) {
    fmt.Println(i, v.(string))
}
// 0 A
// 1 B
// 2 C

// Race — settles with whichever Promise finishes first
fmt.Println("Winner:", vl.Race(p1, p2, p3).GetValue().(string))

// Any — resolves with the first success, or AggregateError if all fail
fmt.Println("Winner:", vl.Any(p1, p2, p3).GetValue().(string))
```

### Async Operations

Wrap the entire Promise chain in a goroutine, not individual callbacks.

```go
done := make(chan struct{})
go func() {
    result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
        time.Sleep(1 * time.Second)
        resolve("async done", nil)
    }).Then(func(value interface{}) (interface{}, error) {
        return value.(string) + "!", nil
    }, nil)

    fmt.Println(result.GetValue()) // async done!
    close(done)
}()
<-done
```

More examples in the [`examples`](./examples) directory.

## License

[MIT](./LICENSE)
