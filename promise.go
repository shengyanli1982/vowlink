package vowlink

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// PromiseState 表示 Promise 的状态
type PromiseState uint8

// Promise 操作的默认处理函数（不可变）
func defaultSuccessHandler(value any) (any, error) { return value, nil }
func defaultErrorHandler(err error) (any, error)   { return nil, err }
func defaultCleanupHandler() error                 { return nil }

// AggregateError 表示错误集合
type AggregateError struct {
	Errors []error
	// 缓存 Error() 的结果，避免每次调用时重新拼接字符串和分配切片。
	// 使用 atomic.Pointer 保证并发调用 Error()/InvalidateError() 时
	// 对缓存的读写不产生数据竞争（竞争下重复构建结果无害）。
	// 注意：Errors 切片自身的并发修改仍需调用方同步；
	// 库内部用法（Any）在结算前完成全部写入，结算后只读。
	cached atomic.Pointer[string]
}

func (ae *AggregateError) Error() string {
	// 如果有缓存且 Errors 未变更，直接返回
	if cached := ae.cached.Load(); cached != nil {
		return *cached
	}

	if len(ae.Errors) == 0 {
		return "All promises were rejected"
	}

	errStrings := make([]string, 0, len(ae.Errors))
	for _, err := range ae.Errors {
		if err == nil {
			errStrings = append(errStrings, "<nil>")
			continue
		}
		errStrings = append(errStrings, err.Error())
	}
	result := "All promises were rejected: " + strings.Join(errStrings, ", ")
	ae.cached.Store(&result)
	return result
}

// InvalidateError 清除 Error() 的缓存，应在 Errors 被修改后调用
func (ae *AggregateError) InvalidateError() {
	ae.cached.Store(nil)
}

func NewAggregateError(capacity int) *AggregateError {
	return &AggregateError{
		Errors: make([]error, 0, capacity),
	}
}

// Promise 状态常量
const (
	Pending   PromiseState = iota // 等待中
	Fulfilled                     // 已完成
	Rejected                      // 已拒绝
)

// String returns the string representation of the PromiseState.
func (s PromiseState) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Fulfilled:
		return "Fulfilled"
	case Rejected:
		return "Rejected"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

// subscriber 注册 Pending 状态 Promise 的回调
type subscriber struct {
	onFulfilled        func(any) (any, error)
	onRejected         func(error) (any, error)
	resolve            func(any, error)
	reject             func(any, error)
	onFulfilledIndexed func(int, any) (any, error)
	onRejectedIndexed  func(int, error) (any, error)
	index              int
}

// noopSettle 是一个空操作的决议函数，用作直接 subscriber 的 resolve/reject 占位符。
// 直接 subscriber 在 onFulfilled/onRejected 回调内部直接处理决议逻辑，
// 因此 resolve/reject 字段仅为满足 dispatchCallback 的接口要求。
func noopSettle(any, error) {}

// Promise 表示一个异步操作
type Promise struct {
	mu          sync.RWMutex
	state       PromiseState
	value       any
	reason      error
	subscribers []subscriber
}

// dispatchSubscriber 分派单个 subscriber，包含 panic 保护。
// 提取为命名函数避免在 settle 循环中重复创建闭包。
func dispatchSubscriber(sub *subscriber, state PromiseState, value any, reason error) {
	defer func() {
		if r := recover(); r != nil {
			sub.reject(nil, fmt.Errorf("subscriber callback panic: %v", r))
		}
	}()
	if sub.onFulfilledIndexed != nil {
		if state == Rejected || reason != nil {
			sub.onRejectedIndexed(sub.index, reason)
		} else {
			sub.onFulfilledIndexed(sub.index, value)
		}
		return
	}
	dispatchCallback(state, value, reason, sub.onFulfilled, sub.onRejected, sub.resolve, sub.reject)
}

// dispatchCallback 根据 Promise 状态选择并调用相应的处理函数，
// 然后根据处理结果决议下游 Promise。
// 此函数提取了 settle() 和 Then() 中重复的回调分发逻辑。
func dispatchCallback(
	state PromiseState, value any, reason error,
	onSuccess func(any) (any, error),
	onError func(error) (any, error),
	resolve func(any, error),
	reject func(any, error),
) {
	var v any
	var e error

	if state == Rejected {
		v, e = onError(reason)
	} else if reason != nil {
		// Fulfilled 但携带错误（resolve(value, err) 模式）
		v, e = onError(reason)
	} else {
		v, e = onSuccess(value)
	}

	if e != nil {
		reject(nil, e)
	} else {
		resolve(v, nil)
	}
}

// dispatchIndexedCallback 根据 Promise 状态调用带索引的回调函数。
// 用于 All/AllSettled 等需要按索引存储结果到共享切片的场景。
// 包含 panic 保护，直接 subscriber 的 panic 会被吞没（不影响外层 Promise）。
func dispatchIndexedCallback(
	index int, state PromiseState, value any, reason error,
	onFulfilled func(int, any) (any, error),
	onRejected func(int, error) (any, error),
) {
	defer func() {
		if r := recover(); r != nil {
			// Panic in direct subscriber callback - swallowed to prevent cascade
			_ = r
		}
	}()
	if state == Rejected {
		onRejected(index, reason)
	} else if reason != nil {
		onRejected(index, reason)
	} else {
		onFulfilled(index, value)
	}
}

// settle 决议 Promise 的状态并同步分发给所有已注册的 subscriber。
//
// 注意：subscriber 回调在调用者的 goroutine 中同步顺序执行。
// 回调函数中不应包含阻塞操作（如死循环、I/O 等待），否则后续
// subscriber 将被阻塞。如需异步处理，请在回调中启动 goroutine。
func (p *Promise) settle(state PromiseState, value any, reason error) {
	p.mu.Lock()
	if p.state != Pending {
		p.mu.Unlock()
		return
	}
	if state == Rejected && reason == nil {
		reason = errors.New("promise rejected")
	}
	p.state = state
	p.value = value
	p.reason = reason
	subs := p.subscribers
	p.subscribers = nil
	p.mu.Unlock()
	for i := range subs {
		dispatchSubscriber(&subs[i], state, value, reason)
	}
}

func (p *Promise) snapshot() (PromiseState, any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.state, p.value, p.reason
}

func (p *Promise) getState() PromiseState {
	state, _, _ := p.snapshot()
	return state
}

// subscribeDirect 直接在 Promise 上注册 subscriber，不创建中间 Promise。
// 适用于 All/AllSettled/Any/Race 等内部聚合操作，可避免 .Then() 创建的
// 中间 Promise 和 executor 闭包的堆分配。
//
// 同步路径（Promise 已 settle）：直接分发回调，无额外分配。
// 异步路径（Promise 仍 Pending）：将回调存入 subscribers 列表，等 settle 时触发。
func (p *Promise) subscribeDirect(
	onFulfilled func(any) (any, error),
	onRejected func(error) (any, error),
) {
	sub := subscriber{
		onFulfilled: onFulfilled,
		onRejected:  onRejected,
		resolve:     noopSettle,
		reject:      noopSettle,
	}

	p.mu.Lock()
	if p.state != Pending {
		s, v, r := p.state, p.value, p.reason
		p.mu.Unlock()
		dispatchSubscriber(&sub, s, v, r)
		return
	}
	p.subscribers = append(p.subscribers, sub)
	p.mu.Unlock()
}

// subscribeDirectIndexed 直接在 Promise 上注册带索引的 subscriber，不创建中间 Promise。
// 与 subscribeDirect 类似，但回调接收额外参数 index，用于 All/AllSettled 等
// 需要按索引存储结果到共享切片的场景。共享回调在循环外创建，仅分配 2 次。
//
// 同步路径：通过 dispatchIndexedCallback 直接调用共享回调（无 per-subscriber 闭包分配）。
// 异步路径：回调与索引直接存入 subscriber 结构（无 per-subscriber 闭包分配）。
func (p *Promise) subscribeDirectIndexed(
	index int,
	onFulfilled func(int, any) (any, error),
	onRejected func(int, error) (any, error),
) {
	p.mu.Lock()
	if p.state != Pending {
		s, v, r := p.state, p.value, p.reason
		p.mu.Unlock()
		dispatchIndexedCallback(index, s, v, r, onFulfilled, onRejected)
		return
	}
	p.subscribers = append(p.subscribers, subscriber{
		onFulfilledIndexed: onFulfilled,
		onRejectedIndexed:  onRejected,
		index:              index,
		resolve:            noopSettle,
		reject:             noopSettle,
	})
	p.mu.Unlock()
}

// NewPromise 使用给定的处理函数创建新的 Promise
//
// 设计说明：
//   - resolve/reject 为直接捕获 p 的内联闭包，每次调用产生 2 次
//     闭包堆分配。这是有意的正确性代价：若复用池化闭包，handler
//     保留 resolve/reject 引用并在结算后再次调用（合法用法）时，
//     会决议被重绑定的其他 Promise，造成静默数据损坏。
//   - 内层 func() 用于 panic recovery，其闭包和 defer 均由编译器
//     在栈上分配（escape analysis 确认 "func literal does not escape"），
//     不产生堆分配开销。
func NewPromise(promiseHandler func(resolve func(any, error), reject func(any, error))) (result *Promise) {
	if promiseHandler == nil {
		return &Promise{
			state:  Rejected,
			reason: errors.New("promise handler cannot be nil"),
		}
	}

	p := &Promise{state: Pending}
	result = p

	resolve := func(value any, reason error) { p.settle(Fulfilled, value, reason) }
	reject := func(value any, reason error) { p.settle(Rejected, value, reason) }

	// 内层 func() + defer 用于隔离 handler panic：
	// 当 handler panic 时，内层 defer recover 后内层函数正常返回，
	// 外层 NewPromise 不受影响，可正常返回 p。
	// escape analysis 确认这两个闭包均在栈上分配，无堆分配。
	func() {
		defer func() {
			if r := recover(); r != nil {
				p.settle(Rejected, nil, fmt.Errorf("promise executor panic: %v", r))
			}
		}()
		promiseHandler(resolve, reject)
	}()

	return p
}

// Then 在 Promise 满足时调用 successHandler，拒绝时调用 errorHandler。
// 若 handler 为 nil 则使用默认透传处理。
// 返回一个新的 Promise，其状态和值由 handler 的返回值决定。
func (p *Promise) Then(successHandler func(any) (any, error), errorHandler func(error) (any, error)) *Promise {
	if successHandler == nil {
		successHandler = defaultSuccessHandler
	}
	if errorHandler == nil {
		errorHandler = defaultErrorHandler
	}

	return NewPromise(func(resolve func(any, error), reject func(any, error)) {
		state, value, reason := p.snapshot()
		switch state {
		case Fulfilled, Rejected:
			dispatchCallback(state, value, reason, successHandler, errorHandler, resolve, reject)
		case Pending:
			p.mu.Lock()
			if p.state != Pending {
				s, v, r := p.state, p.value, p.reason
				p.mu.Unlock()
				dispatchCallback(s, v, r, successHandler, errorHandler, resolve, reject)
			} else {
				p.subscribers = append(p.subscribers, subscriber{
					onFulfilled: successHandler,
					onRejected:  errorHandler,
					resolve:     resolve,
					reject:      reject,
				})
				p.mu.Unlock()
			}
		}
	})
}

// Catch 在 Promise 被拒绝时调用 errorHandler，返回可恢复的新 Promise。
func (p *Promise) Catch(errorHandler func(error) (any, error)) *Promise {
	return p.Then(nil, errorHandler)
}

// Finally 无论 Promise 满足或拒绝都会调用 cleanupHandler。
// 若 cleanupHandler 返回 error，则下游 Promise 以该错误拒绝。
func (p *Promise) Finally(cleanupHandler func() error) *Promise {
	if cleanupHandler == nil {
		cleanupHandler = defaultCleanupHandler
	}

	return p.Then(
		func(value any) (any, error) {
			err := cleanupHandler()
			if err != nil {
				return nil, err
			}
			return value, nil
		},
		func(reason error) (any, error) {
			err := cleanupHandler()
			if err != nil {
				return nil, errors.Join(reason, err)
			}
			return nil, reason
		},
	)
}

// GetValue 返回 Promise 的满足值。若 Promise 尚未完成则返回 nil（线程安全）。
func (p *Promise) GetValue() any {
	_, value, _ := p.snapshot()
	return value
}

// GetReason 返回 Promise 的拒绝原因。若 Promise 尚未完成或为满足则返回 nil（线程安全）。
func (p *Promise) GetReason() error {
	_, _, reason := p.snapshot()
	return reason
}

// countNonNil 原地统计非 nil Promise 的数量，不分配新切片
func countNonNil(promises []*Promise) int {
	count := 0
	for _, p := range promises {
		if p != nil {
			count++
		}
	}
	return count
}

// All 等待所有 Promise 完成
// 如果任何一个 Promise 被拒绝，结果 Promise 也会被拒绝
//
// 性能优化：
//   - 使用 countNonNil 原地遍历，避免 filterNilPromises 的切片分配
//   - 使用 subscribeDirectIndexed 消除中间 Promise 和 executor 闭包分配
//   - 共享回调在循环外创建，同步路径仅需 2 次闭包分配（vs N×3 次）
func All(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(any, error), reject func(any, error)) {
		count := countNonNil(promises)
		if count == 0 {
			resolve([]any{}, nil)
			return
		}

		values := make([]any, count)
		var pending atomic.Int32
		var rejected atomic.Bool
		pending.Store(int32(count))

		// 内存序：每个回调先写 values[i]，再对同一原子计数器执行 Add(-1)（release 语义）；
		// 观察到 Add 返回 0 的操作在顺序一致性总序中位于所有先前 Add 之后（acquire 语义），
		// 因此可见全部 slot 写入；settle 幂等兜底。

		// 共享回调在循环外创建（仅 2 次闭包分配，vs 循环内 N×2 次）
		onFulfilled := func(i int, value any) (any, error) {
			values[i] = value
			if pending.Add(-1) == 0 {
				resolve(values, nil)
			}
			return nil, nil
		}
		onRejected := func(_ int, reason error) (any, error) {
			if rejected.CompareAndSwap(false, true) {
				reject(nil, reason)
			}
			return nil, nil
		}

		nonNilIdx := 0
		for _, promise := range promises {
			if promise == nil {
				continue
			}
			promise.subscribeDirectIndexed(nonNilIdx, onFulfilled, onRejected)
			nonNilIdx++
		}
	})
}

// AllSettled 等待所有 Promise 完成，无论其状态如何
//
// 性能优化：
//   - 使用 countNonNil 原地遍历，避免 filterNilPromises 的切片分配
//   - 使用 subscribeDirectIndexed 消除中间 Promise 和 executor 闭包分配
//   - 共享回调在循环外创建，同步路径仅需 2 次闭包分配（vs N×2 次）
func AllSettled(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(any, error), reject func(any, error)) {
		count := countNonNil(promises)
		if count == 0 {
			resolve([]any{}, nil)
			return
		}

		values := make([]any, count)
		var pending atomic.Int32
		pending.Store(int32(count))

		// 内存序：每个回调先写 values[i]，再对同一原子计数器执行 Add(-1)（release 语义）；
		// 观察到 Add 返回 0 的操作在顺序一致性总序中位于所有先前 Add 之后（acquire 语义），
		// 因此可见全部 slot 写入；settle 幂等兜底。

		// 共享回调在循环外创建（仅 2 次闭包分配）
		onFulfilled := func(i int, value any) (any, error) {
			values[i] = value
			if pending.Add(-1) == 0 {
				resolve(values, nil)
			}
			return nil, nil
		}
		onRejected := func(i int, reason error) (any, error) {
			values[i] = reason
			if pending.Add(-1) == 0 {
				resolve(values, nil)
			}
			return nil, nil
		}

		nonNilIdx := 0
		for _, promise := range promises {
			if promise == nil {
				continue
			}
			promise.subscribeDirectIndexed(nonNilIdx, onFulfilled, onRejected)
			nonNilIdx++
		}
	})
}

// Any 返回一个在任意输入 Promise 成功时完成的 Promise
// 如果所有 Promise 都被拒绝，返回一个 AggregateError
//
// 性能优化：
//   - 使用 countNonNil 原地遍历，避免 filterNilPromises 的切片分配
//   - 使用 subscribeDirectIndexed 消除中间 Promise 和 executor 闭包分配
//   - 共享回调在循环外创建，同步路径仅需 2 次闭包分配
func Any(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(any, error), reject func(any, error)) {
		count := countNonNil(promises)
		if count == 0 {
			reject(nil, NewAggregateError(0))
			return
		}

		aggErr := NewAggregateError(count)
		aggErr.Errors = aggErr.Errors[:count]
		var pending atomic.Int32
		var done atomic.Bool
		pending.Store(int32(count))

		// 共享回调在循环外创建（索引经参数传入，所有 promise 共享）
		onFulfilled := func(_ int, value any) (any, error) {
			if done.CompareAndSwap(false, true) {
				resolve(value, nil)
			}
			return nil, nil
		}
		onRejected := func(i int, reason error) (any, error) {
			aggErr.Errors[i] = reason
			// 不变式：每个输入 Promise 恰好分派一次回调且仅走一个分支，pending.Add(-1)==0 ⟺ 全部输入均走 onRejected（含 resolve(value, err) 转拒分支，该分支同样只递减计数而不触碰 done）⟹ onFulfilled 从未运行 ⟹ done 必为 false，故此 CAS 可证冗余。
			if pending.Add(-1) == 0 {
				reject(nil, aggErr)
			}
			return nil, nil
		}

		nonNilIdx := 0
		for _, promise := range promises {
			if promise == nil {
				continue
			}
			promise.subscribeDirectIndexed(nonNilIdx, onFulfilled, onRejected)
			nonNilIdx++
		}
	})
}

// Race 返回一个与第一个完成的 Promise 具有相同状态的 Promise
//
// 性能优化：
//   - 使用 countNonNil 原地遍历，避免 filterNilPromises 的切片分配
//   - 使用 subscribeDirect 消除中间 Promise 和 executor 闭包分配
//   - 共享回调在循环外创建（无需索引），同步路径仅需 2 次闭包分配
func Race(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(any, error), reject func(any, error)) {
		count := countNonNil(promises)
		if count == 0 {
			resolve(nil, nil)
			return
		}

		var done atomic.Bool

		// 共享回调在循环外创建（无索引依赖，所有 promise 共享）
		onFulfilled := func(value any) (any, error) {
			if done.CompareAndSwap(false, true) {
				resolve(value, nil)
			}
			return nil, nil
		}
		onRejected := func(reason error) (any, error) {
			if done.CompareAndSwap(false, true) {
				reject(nil, reason)
			}
			return nil, nil
		}

		for _, promise := range promises {
			if promise == nil {
				continue
			}
			promise.subscribeDirect(onFulfilled, onRejected)
		}
	})
}
