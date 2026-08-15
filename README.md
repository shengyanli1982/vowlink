<div align="center">
    <img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/vowlink.svg)](https://pkg.go.dev/github.com/shengyanli1982/vowlink)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shengyanli1982/vowlink)

A lightweight Promise library for Go, inspired by the JavaScript Promise API. VowLink replaces callback hell with a fluent, chainable API for both synchronous and asynchronous workflows.

**Features:**

- **JS-familiar API** — `Then` / `Catch` / `Finally` / `All` / `Race` / `Any` / `AllSettled`
- **Zero dependencies** — standard library only
- **Thread-safe** — `sync.RWMutex` protects all state transitions
- **Immutable state** — once settled (`Fulfilled` / `Rejected`), state never changes
- **High performance** — lock-free atomic dispatch in collection combinators, constant per-link allocations in chained APIs

**Requires Go 1.23+.**

## Install

```bash
go get github.com/shengyanli1982/vowlink
```

## Quick Start

```go
import vl "github.com/shengyanli1982/vowlink"

result := vl.NewPromise(func(resolve, reject func(any, error)) {
    resolve("hello", nil)
}).Then(func(v any) (any, error) {
    return v.(string) + " vowlink", nil
}, nil)

fmt.Println(result.GetValue()) // hello vowlink
```

## How It Works

### State Machine

```
         resolve()                reject()
            │                        │
    ┌───────▼────────┐      ┌───────▼────────┐
    │   Fulfilled    │      │    Rejected     │
    └────────────────┘      └────────────────┘
                 ▲                    ▲
                 └────── Pending ─────┘
```

Every Promise starts as `Pending` and transitions to exactly one terminal state. Once settled, the state is **permanent** and thread-safe.

### Handler Signatures

All handlers follow a unified signature for maximum flexibility:

```go
func onFulfilled(value any) (any, error)
func onRejected(reason error) (any, error)
```

- Return `(value, nil)` → downstream **Fulfills** with `value`
- Return `(nil, err)` → downstream **Rejects** with `err`
- In `Catch`, returning a value with `nil` error **recovers** the chain back to Fulfilled

### Resolve & Reject

Both `resolve` and `reject` accept `(value, error)`. The dispatch logic determines the path:

| Call                  | Effect                                                     |
| --------------------- | ---------------------------------------------------------- |
| `resolve(value, nil)` | Fulfills; `Then` calls `onFulfilled`                       |
| `resolve(value, err)` | Fulfills; `Then` calls `onRejected` (error takes priority) |
| `reject(nil, err)`    | Rejects; `Then` calls `onRejected`                         |

Inside combinators, `resolve(value, err)` is treated as a rejection: `All`/`Race` reject with `err`, `AllSettled` stores `err` in its result slice, and `Any` records `err` in the `AggregateError`. As in `Then`, the error takes priority and `value` is discarded.

## API Reference

### Instance Methods

| Method                          | Description                                                          |
| ------------------------------- | -------------------------------------------------------------------- |
| `Then(onFulfilled, onRejected)` | Register callbacks; returns new Promise. Either handler may be `nil` |
| `Catch(onRejected)`             | Shorthand for `Then(nil, onRejected)`                                |
| `Finally(cleanup)`              | Always runs cleanup handler; returning `error` rejects downstream    |
| `GetValue() any`        | Returns resolved value, or `nil` if not yet fulfilled                |
| `GetReason() error`             | Returns rejection reason, or `nil` if not rejected                   |

### Combinators

| Function                  | Settles When                   | Result                                    |
| ------------------------- | ------------------------------ | ----------------------------------------- |
| `All(p1, p2, ...)`        | All fulfill OR first rejects   | `[]any` of values, or first error |
| `Race(p1, p2, ...)`       | First settles (any state)      | Value/error of the winner                 |
| `Any(p1, p2, ...)`        | First fulfills OR all reject   | First value, or `*AggregateError`         |
| `AllSettled(p1, p2, ...)` | All settle (fulfill or reject) | `[]any` mixing values and errors  |

**Empty input:** `Race()` fulfills immediately with `nil` — a deliberate design choice, unlike JavaScript where an empty race stays pending forever. `Any()` rejects with an empty `*AggregateError`.

## Usage Examples

### Chain & Error Handling

```go
result := vl.NewPromise(func(resolve, reject func(any, error)) {
    reject(nil, fmt.Errorf("network timeout"))
}).Then(func(v any) (any, error) {
    return v, nil // skipped on rejection
}, nil).Catch(func(err error) (any, error) {
    return nil, fmt.Errorf("wrapped: %w", err)
})

fmt.Println(result.GetReason()) // wrapped: network timeout
```

### Error Recovery

Returning a value with `nil` error in `Catch` **recovers** the chain:

```go
result := vl.NewPromise(func(resolve, reject func(any, error)) {
    reject(nil, fmt.Errorf("disk full"))
}).Catch(func(err error) (any, error) {
    return "fallback data", nil // recover: chain becomes Fulfilled
}).Then(func(v any) (any, error) {
    return fmt.Sprintf("got: %v", v), nil
}, nil)

fmt.Println(result.GetValue()) // got: fallback data
```

### Cleanup with Finally

`Finally` runs regardless of state. If the cleanup handler returns an error, downstream rejects:

```go
vl.NewPromise(func(resolve, reject func(any, error)) {
    resolve("data", nil)
}).Finally(func() error {
    fmt.Println("cleanup done") // always prints
    return nil
}).Then(func(v any) (any, error) {
    fmt.Println(v) // data
    return v, nil
}, nil)
```

### Async Operations

Use goroutines inside the executor for async work. VowLink dispatches subscribers when settled:

```go
p := vl.NewPromise(func(resolve, reject func(any, error)) {
    go func() {
        time.Sleep(100 * time.Millisecond)
        resolve("async result", nil)
    }()
})

p.Then(func(v any) (any, error) {
    fmt.Println(v) // async result (called when settled)
    return v, nil
}, nil)
```

### Collection Combinators

```go
p1 := vl.NewPromise(func(r, _ func(any, error)) { r("A", nil) })
p2 := vl.NewPromise(func(r, _ func(any, error)) { r("B", nil) })
p3 := vl.NewPromise(func(r, _ func(any, error)) { r("C", nil) })

// All: collect results in order
values := vl.All(p1, p2, p3).GetValue().([]any)
// [A B C]

// Race: first settled wins
winner := vl.Race(p1, p2, p3).GetValue()
// A (or whichever finishes first for async promises)

// Any: first fulfilled, or *AggregateError if all rejected
result := vl.Any(p1, p2, p3)

// AllSettled: wait for all, collect values and errors
mixed := vl.AllSettled(p1, p2, p3).GetValue().([]any)
```

## Concurrency

VowLink is designed for concurrent use. Multiple goroutines can safely call `resolve`/`reject` concurrently — only the first one wins, and all subscribers are dispatched.

**Guidelines:**

- Subscriber callbacks in `Then`/`Catch`/`Finally` execute **synchronously** on the caller's goroutine
- Do not block inside callbacks — spawn a goroutine if you need async work
- `GetValue()` and `GetReason()` are safe to call from any goroutine at any time
- **Early-settle retention** — after `All`/`Any`/`Race` settles early, subscribers on the remaining pending inputs stay attached until each input settles (no unsubscription, matching JavaScript semantics)
- **Settle cascade depth** — callbacks run synchronously on the settling goroutine, so settling a pending `Then` chain recurses to a depth equal to the chain length; Go's growable stacks support real-world depths (tests cover 10000 links), but consider splitting extremely deep chains or dispatching them asynchronously

## Benchmarks

```
$ go test -bench=. -benchmem -count=5
BenchmarkPromiseThenChain/chain=1          6 allocs/op    256 B/op
BenchmarkPromiseThenChain/chain=32        99 allocs/op   4.1 KB/op
BenchmarkPromiseFinally                    8 allocs/op    288 B/op
BenchmarkPromiseAll/size=128               8 allocs/op   2.5 KB/op
BenchmarkPromiseAllSettled/size=128        8 allocs/op   2.5 KB/op
BenchmarkPromiseAny/first-fulfilled        8 allocs/op   0.7 KB/op
BenchmarkPromiseRace/size=128              6 allocs/op   0.2 KB/op
BenchmarkAggregateError/cached/size=64     0 allocs/op      0 B/op
BenchmarkPromiseThenPending                7 allocs/op    320 B/op
BenchmarkPromiseAllPending/size=128      522 allocs/op  28.7 KB/op
```

NewPromise deliberately spends 2 closure allocations per call for correctness, preventing cross-Promise contamination from pooled closures; the collection combinators are lock-free atomic, and AllPending measures the async subscription path.

## Examples

See the [`examples/`](./examples/) directory for 13 runnable demos covering all features.

## License

[MIT](./LICENSE)
