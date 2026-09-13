package vowlink

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

var benchmarkErr = errors.New("benchmark rejection")

func benchmarkIncrement(value any) (any, error) {
	return value.(int) + 1, nil
}

func benchmarkPropagateError(err error) (any, error) {
	return nil, err
}

func benchmarkCleanupNoop() error {
	return nil
}

func makeFulfilledPromises(size int) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve(i, nil)
		})
	}
	return promises
}

func makeRejectedPromises(size int, err error) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, err)
		})
	}
	return promises
}

func makeMixedPromises(size int, err error) []*Promise {
	promises := make([]*Promise, size)
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
				resolve(i, nil)
			})
			continue
		}
		promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
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

	promises[0] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
		resolve("winner", nil)
	})

	for i := 1; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, err)
		})
	}

	return promises
}

// makePendingPromises 创建 size 个 pending 状态的 promise，并返回它们的 resolve/reject 句柄。
// executor 内不结算，由调用方在订阅完成后按需结算，用于测量 pending 订阅 → settle 批量分发的真实异步路径。
func makePendingPromises(size int) ([]*Promise, []func(any, error), []func(any, error)) {
	promises := make([]*Promise, size)
	resolvers := make([]func(any, error), size)
	rejecters := make([]func(any, error), size)
	for i := 0; i < size; i++ {
		promises[i] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolvers[i] = resolve
			rejecters[i] = reject
		})
	}
	return promises, resolvers, rejecters
}

func BenchmarkPromiseThenChain(b *testing.B) {
	for _, chainLength := range []int{1, 8, 32} {
		b.Run(fmt.Sprintf("chain=%d", chainLength), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		b.Run(fmt.Sprintf("chain=%d", chainLength), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("ok", nil)
		}).Finally(benchmarkCleanupNoop)

		if i == 0 && (p.GetReason() != nil || p.GetValue() != "ok") {
			b.Fatalf("unexpected Finally result: value=%v reason=%v", p.GetValue(), p.GetReason())
		}
	}
}

func BenchmarkPromiseAll(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		promises := makeFulfilledPromises(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := All(promises...)
				if i == 0 {
					values, ok := p.GetValue().([]any)
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
		promises := makeMixedPromises(size, benchmarkErr)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := AllSettled(promises...)
				if i == 0 {
					values, ok := p.GetValue().([]any)
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

func BenchmarkAggregateError(b *testing.B) {
	for _, size := range []int{1, 8, 64} {
		b.Run(fmt.Sprintf("build/size=%d", size), func(b *testing.B) {
			// 成员错误在计时循环外预生成：fmt 构造不进入测量，
			// build 变体只测 NewAggregateError + Error()（strings.Join）库路径。
			members := make([]error, size)
			for j := 0; j < size; j++ {
				members[j] = fmt.Errorf("error-%d", j)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				aggregateErr := NewAggregateError(size)
				aggregateErr.Errors = append(aggregateErr.Errors, members...)
				message := aggregateErr.Error()
				if i == 0 && !strings.Contains(message, "error-0") {
					b.Fatalf("unexpected AggregateError build message: %q", message)
				}
			}
		})

		b.Run(fmt.Sprintf("cached/size=%d", size), func(b *testing.B) {
			aggregateErr := NewAggregateError(size)
			for j := 0; j < size; j++ {
				aggregateErr.Errors = append(aggregateErr.Errors, errors.New(fmt.Sprintf("error-%d", j)))
			}
			_ = aggregateErr.Error()
			b.ReportAllocs()
			b.ResetTimer()
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

		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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

func BenchmarkPromiseThenPending(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var asyncResolve func(any, error)
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			asyncResolve = resolve
		})
		result := p.Then(benchmarkIncrement, nil)
		asyncResolve(0, nil)

		if i == 0 {
			finalValue, ok := result.GetValue().(int)
			if !ok || finalValue != 1 || result.GetReason() != nil {
				b.Fatalf("unexpected pending Then result: value=%v reason=%v", result.GetValue(), result.GetReason())
			}
		}
	}
}

func BenchmarkPromiseAllPending(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				promises := make([]*Promise, size)
				resolvers := make([]func(any, error), size)
				for j := 0; j < size; j++ {
					promises[j] = NewPromise(func(resolve func(any, error), reject func(any, error)) {
						resolvers[j] = resolve
					})
				}

				p := All(promises...)
				for j := 0; j < size; j++ {
					resolvers[j](j, nil)
				}

				if i == 0 {
					values, ok := p.GetValue().([]any)
					if !ok || p.GetReason() != nil || len(values) != size {
						b.Fatalf("unexpected pending All result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
					for j := 0; j < size; j++ {
						if values[j] != j {
							b.Fatalf("unexpected pending All value at index %d: %v", j, values[j])
						}
					}
				}
			}
		})
	}
}

// BenchmarkPromiseRacePending 测量 Race 的 pending 路径：计时循环内创建 pending promise，
// Race 订阅存储 subscriber 后结算首个输入，触发 done.CAS 选主与存储分发。
// 与 BenchmarkPromiseRace（预结算输入、热同步分发）互补；创建在循环内是有意的，
// 覆盖"创建 + pending 订阅 + 结算分发"的真实异步主成本（PERF-REPORT §1.4）。
func BenchmarkPromiseRacePending(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				promises, resolvers, _ := makePendingPromises(size)
				p := Race(promises...)
				resolvers[0]("winner", nil)

				if i == 0 && (p.GetValue() != "winner" || p.GetReason() != nil) {
					b.Fatalf("unexpected RacePending result: value=%v reason=%v", p.GetValue(), p.GetReason())
				}
			}
		})
	}
}

// BenchmarkPromiseAnyPending 测量 Any 的 pending 路径：首个 fulfill 获胜
// （done.CAS + resolve 短路），其余输入保持 pending。
func BenchmarkPromiseAnyPending(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				promises, resolvers, _ := makePendingPromises(size)
				p := Any(promises...)
				resolvers[0]("first", nil)

				if i == 0 && (p.GetValue() != "first" || p.GetReason() != nil) {
					b.Fatalf("unexpected AnyPending result: value=%v reason=%v", p.GetValue(), p.GetReason())
				}
			}
		})
	}
}

// BenchmarkPromiseAllSettledPending 测量 AllSettled 的 pending 路径：
// 全部输入逐个结算，pending 计数器归零后 resolve 聚合结果。
func BenchmarkPromiseAllSettledPending(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				promises, resolvers, _ := makePendingPromises(size)
				p := AllSettled(promises...)
				for j := 0; j < size; j++ {
					resolvers[j](j, nil)
				}

				if i == 0 {
					values, ok := p.GetValue().([]any)
					if !ok || p.GetReason() != nil || len(values) != size {
						b.Fatalf("unexpected AllSettledPending result: value=%v reason=%v", p.GetValue(), p.GetReason())
					}
					for j := 0; j < size; j++ {
						if values[j] != j {
							b.Fatalf("unexpected AllSettledPending value at index %d: %v", j, values[j])
						}
					}
				}
			}
		})
	}
}

// BenchmarkPromiseAllRejectSettled 测量 All 的拒绝短路路径（预结算输入）：
// 全部输入已 rejected，订阅时同步走 onRejected（rejected.CAS 仅首个成功 + reject），
// 这是 All 的 fail-fast 定义性行为，现有 BenchmarkPromiseAll（全成功）从不触发。
func BenchmarkPromiseAllRejectSettled(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		promises := makeRejectedPromises(size, benchmarkErr)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := All(promises...)
				if i == 0 && !errors.Is(p.GetReason(), benchmarkErr) {
					b.Fatalf("unexpected AllRejectSettled result: reason=%v", p.GetReason())
				}
			}
		})
	}
}

// BenchmarkPromiseAllRejectPending 测量 All 的拒绝短路路径（pending 输入）：
// 前一半输入先结算成功，中点输入 reject 触发短路，其余保持 pending
// （结算时序与 PERF-REPORT §3.2 基线一致，保证数字可比）。
func BenchmarkPromiseAllRejectPending(b *testing.B) {
	for _, size := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				promises, resolvers, rejecters := makePendingPromises(size)
				p := All(promises...)
				for j := 0; j < size/2; j++ {
					resolvers[j](j, nil)
				}
				rejecters[size/2](nil, benchmarkErr)

				if i == 0 && !errors.Is(p.GetReason(), benchmarkErr) {
					b.Fatalf("unexpected AllRejectPending result: reason=%v", p.GetReason())
				}
			}
		})
	}
}

// BenchmarkPromiseThenFanoutPending 测量 Then 的多 subscriber pending fanout：
// 单个 pending promise 挂 fanout 个 Then（subscribers 切片增长），
// 一次 settle 批量 drain 分发给全部下游。现有 BenchmarkPromiseThenPending 只测单个 Then。
func BenchmarkPromiseThenFanoutPending(b *testing.B) {
	for _, fanout := range []int{4, 32, 128} {
		b.Run(fmt.Sprintf("fanout=%d", fanout), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var asyncResolve func(any, error)
				p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
					asyncResolve = resolve
				})
				results := make([]*Promise, fanout)
				for j := 0; j < fanout; j++ {
					results[j] = p.Then(benchmarkIncrement, nil)
				}
				asyncResolve(0, nil)

				if i == 0 {
					for j := 0; j < fanout; j++ {
						value, ok := results[j].GetValue().(int)
						if !ok || value != 1 || results[j].GetReason() != nil {
							b.Fatalf("unexpected ThenFanoutPending result at %d: value=%v reason=%v",
								j, results[j].GetValue(), results[j].GetReason())
						}
					}
				}
			}
		})
	}
}
