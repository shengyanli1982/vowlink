package vowlink

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

var benchmarkErr = errors.New("benchmark rejection")

func benchmarkIncrement(value interface{}) (interface{}, error) {
	return value.(int) + 1, nil
}

func benchmarkPropagateError(err error) (interface{}, error) {
	return nil, err
}

func benchmarkCleanupNoop() error {
	return nil
}

func makeFulfilledPromises(size int) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		value := i
		promises[i] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve(value, nil)
		})
	}
	return promises
}

func makeRejectedPromises(size int, err error) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, err)
		})
	}
	return promises
}

func makeMixedPromises(size int, err error) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			value := i
			promises[i] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
				resolve(value, nil)
			})
			continue
		}
		promises[i] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, err)
		})
	}
	return promises
}

func makeRacePromises(size int, err error) []*Promise {
	promises := make([]*Promise, size)
	if size == 0 {
		return promises
	}

	promises[0] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		resolve("winner", nil)
	})

	for i := 1; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, err)
		})
	}

	return promises
}

func BenchmarkPromiseThenChain(b *testing.B) {
	for _, chainLength := range []int{1, 8, 32} {
		chainLength := chainLength
		b.Run(fmt.Sprintf("chain=%d", chainLength), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
					resolve(0, nil)
				})
				for j := 0; j < chainLength; j++ {
					p = p.Then(benchmarkIncrement, nil)
				}

				if i == 0 {
					finalValue, ok := p.GetValue().(int)
					if !ok || finalValue != chainLength || p.GetReason() != nil {
						b.Fatalf("unexpected Then chain result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
				}
			}
		})
	}
}

func BenchmarkPromiseCatchChain(b *testing.B) {
	for _, chainLength := range []int{1, 8, 32} {
		chainLength := chainLength
		b.Run(fmt.Sprintf("chain=%d", chainLength), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
					reject(nil, benchmarkErr)
				})
				for j := 0; j < chainLength; j++ {
					p = p.Catch(benchmarkPropagateError)
				}

				if i == 0 {
					if p.GetReason() == nil || p.GetValue() != nil {
						b.Fatalf("unexpected Catch chain result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
				}
			}
		})
	}
}

func BenchmarkPromiseFinally(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("ok", nil)
		}).Finally(benchmarkCleanupNoop)

		if i == 0 && (p.GetReason() != nil || p.GetValue() != "ok") {
			b.Fatalf("unexpected Finally result: value=%v reason=%v", p.GetValue(), p.GetReason())
		}
	}
}

func BenchmarkPromiseAll(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		size := size
		promises := makeFulfilledPromises(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := All(promises...)
				if i == 0 {
					values, ok := p.GetValue().([]interface{})
					if !ok || p.GetReason() != nil || len(values) != size {
						b.Fatalf("unexpected All result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
				}
			}
		})
	}
}

func BenchmarkPromiseAllSettled(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		size := size
		promises := makeMixedPromises(size, benchmarkErr)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := AllSettled(promises...)
				if i == 0 {
					values, ok := p.GetValue().([]interface{})
					if !ok || p.GetReason() != nil || len(values) != size {
						b.Fatalf("unexpected AllSettled result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
				}
			}
		})
	}
}

func BenchmarkPromiseAny(b *testing.B) {
	b.Run("first-fulfilled", func(b *testing.B) {
		promises := append(makeFulfilledPromises(1), makeRejectedPromises(31, benchmarkErr)...)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p := Any(promises...)
			if i == 0 && (p.GetReason() != nil || p.GetValue() == nil) {
				b.Fatalf("unexpected Any result: value=%v reason=%v", p.GetValue(), p.GetReason())
			}
		}
	})

	b.Run("all-rejected", func(b *testing.B) {
		promises := makeRejectedPromises(32, benchmarkErr)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p := Any(promises...)
			if i == 0 {
				if p.GetReason() == nil {
					b.Fatalf("unexpected Any all-rejected result: reason is nil")
				}
				if _, ok := p.GetReason().(*AggregateError); !ok {
					b.Fatalf("expected AggregateError, got %T", p.GetReason())
				}
			}
		}
	})
}

func BenchmarkPromiseRace(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		size := size
		promises := makeRacePromises(size, benchmarkErr)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := Race(promises...)
				if i == 0 && (p.GetReason() != nil || p.GetValue() != "winner") {
					b.Fatalf("unexpected Race result: value=%v reason=%v", p.GetValue(), p.GetReason())
				}
			}
		})
	}
}

func BenchmarkAggregateErrorError(b *testing.B) {
	for _, size := range []int{1, 8, 64} {
		size := size
		aggregateErr := NewAggregateError(size)
		for i := 0; i < size; i++ {
			aggregateErr.Errors = append(aggregateErr.Errors, errors.New(fmt.Sprintf("error-%d", i)))
		}

		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = aggregateErr.Error()
			}
		})
	}
}

func BenchmarkPromiseConcurrentSettle(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(2)

		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			go func() {
				defer wg.Done()
				resolve("ok", nil)
			}()
			go func() {
				defer wg.Done()
				reject(nil, benchmarkErr)
			}()
		})

		wg.Wait()
		_ = p.getState()
		_ = p.GetValue()
		_ = p.GetReason()
	}
}
