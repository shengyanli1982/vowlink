package vowlink

import (
	"strings"
	"sync"
)

// PromiseState 表示 Promise 的状态
type PromiseState uint8

// Promise 操作的默认处理函数
var (
	defaultSuccessHandler = func(value interface{}) (interface{}, error) { return value, nil }
	defaultErrorHandler   = func(err error) (interface{}, error) { return nil, err }
	defaultCleanupHandler = func() error { return nil }
)

// AggregateError 表示错误集合
type AggregateError struct {
	Errors []error
}

func (ae *AggregateError) Error() string {
	if len(ae.Errors) == 0 {
		return "All promises were rejected"
	}

	errStrings := make([]string, 0, len(ae.Errors))
	for _, err := range ae.Errors {
		if err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}
	return "All promises were rejected: " + strings.Join(errStrings, ", ")
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

// subscriber 注册 Pending 状态 Promise 的回调
type subscriber struct {
	onFulfilled func(interface{}) (interface{}, error)
	onRejected  func(error) (interface{}, error)
	resolve     func(interface{}, error)
	reject      func(interface{}, error)
}

// Promise 表示一个异步操作
type Promise struct {
	mu          sync.RWMutex
	state       PromiseState
	value       interface{}
	reason      error
	subscribers []subscriber
}

func (p *Promise) settle(state PromiseState, value interface{}, reason error) {
	p.mu.Lock()
	if p.state != Pending {
		p.mu.Unlock()
		return
	}
	p.state = state
	p.value = value
	p.reason = reason
	subs := p.subscribers
	p.subscribers = nil
	p.mu.Unlock()
	for _, sub := range subs {
		if reason != nil {
			sub.reject(sub.onRejected(reason))
		} else {
			sub.resolve(sub.onFulfilled(value))
		}
	}
}

func (p *Promise) snapshot() (PromiseState, interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.state, p.value, p.reason
}

func (p *Promise) getState() PromiseState {
	state, _, _ := p.snapshot()
	return state
}

func (p *Promise) resolve(value interface{}, reason error) {
	p.settle(Fulfilled, value, reason)
}

func (p *Promise) reject(value interface{}, reason error) {
	p.settle(Rejected, value, reason)
}

// NewPromise 使用给定的处理函数创建新的 Promise
func NewPromise(promiseHandler func(resolve func(interface{}, error), reject func(interface{}, error))) *Promise {
	if promiseHandler == nil {
		return nil
	}

	p := &Promise{state: Pending}

	promiseHandler(p.resolve, p.reject)

	return p
}

// Then 注册 Promise 完成时要调用的回调函数
func (p *Promise) Then(successHandler func(interface{}) (interface{}, error), errorHandler func(error) (interface{}, error)) *Promise {
	if successHandler == nil {
		successHandler = defaultSuccessHandler
	}
	if errorHandler == nil {
		errorHandler = defaultErrorHandler
	}

	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		state, value, reason := p.snapshot()
		switch state {
		case Fulfilled, Rejected:
			if reason != nil {
				reject(errorHandler(reason))
			} else {
				resolve(successHandler(value))
			}
		case Pending:
			p.mu.Lock()
			if p.state != Pending {
				value, reason := p.value, p.reason
				p.mu.Unlock()
				if reason != nil {
					reject(errorHandler(reason))
				} else {
					resolve(successHandler(value))
				}
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

// Catch 注册 Promise 被拒绝时要调用的回调函数
func (p *Promise) Catch(errorHandler func(error) (interface{}, error)) *Promise {
	return p.Then(nil, errorHandler)
}

// Finally 注册无论 Promise 状态如何都会调用的清理回调函数
func (p *Promise) Finally(cleanupHandler func() error) *Promise {
	if cleanupHandler == nil {
		cleanupHandler = defaultCleanupHandler
	}

	return p.Then(
		func(value interface{}) (interface{}, error) {
			err := cleanupHandler()
			if err != nil {
				return nil, err
			}
			return value, nil
		},
		func(reason error) (interface{}, error) {
			err := cleanupHandler()
			if err != nil {
				return nil, err
			}
			return nil, reason
		},
	)
}

func (p *Promise) GetValue() interface{} {
	_, value, _ := p.snapshot()
	return value
}

func (p *Promise) GetReason() error {
	_, _, reason := p.snapshot()
	return reason
}

// All 等待所有 Promise 完成
// 如果任何一个 Promise 被拒绝，结果 Promise 也会被拒绝
func All(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		filtered := make([]*Promise, 0, len(promises))
		for _, p := range promises {
			if p != nil {
				filtered = append(filtered, p)
			}
		}
		promises = filtered

		if len(promises) == 0 {
			resolve([]interface{}{}, nil)
			return
		}

		values := make([]interface{}, len(promises))
		pendingCount := len(promises)
		isCompleted := false
		var mu sync.Mutex

		for i, promise := range promises {
			promise.Then(func(value interface{}) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				values[i] = value
				pendingCount--
				if pendingCount == 0 {
					isCompleted = true
					resolve(values, nil)
				}
				return nil, nil
			}, func(reason error) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				isCompleted = true
				reject(nil, reason)
				return nil, nil
			})
		}
	})
}

// AllSettled 等待所有 Promise 完成，无论其状态如何
func AllSettled(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		filtered := make([]*Promise, 0, len(promises))
		for _, p := range promises {
			if p != nil {
				filtered = append(filtered, p)
			}
		}
		promises = filtered

		if len(promises) == 0 {
			resolve([]interface{}{}, nil)
			return
		}

		values := make([]interface{}, len(promises))
		pendingCount := len(promises)
		var mu sync.Mutex

		for i, promise := range promises {
			promise.Then(func(value interface{}) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				values[i] = value
				pendingCount--
				if pendingCount == 0 {
					resolve(values, nil)
				}
				return nil, nil
			}, func(reason error) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				values[i] = reason
				pendingCount--
				if pendingCount == 0 {
					resolve(values, nil)
				}
				return nil, nil
			})
		}
	})
}

// Any 返回一个在任意输入 Promise 成功时完成的 Promise
// 如果所有 Promise 都被拒绝，返回一个 AggregateError
func Any(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		filtered := make([]*Promise, 0, len(promises))
		for _, p := range promises {
			if p != nil {
				filtered = append(filtered, p)
			}
		}
		promises = filtered

		if len(promises) == 0 {
			reject(nil, NewAggregateError(0))
			return
		}

		errors := NewAggregateError(len(promises))
		pendingCount := len(promises)
		isCompleted := false
		var mu sync.Mutex

		for _, promise := range promises {
			promise.Then(func(value interface{}) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				isCompleted = true
				resolve(value, nil)
				return nil, nil
			}, func(reason error) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				errors.Errors = append(errors.Errors, reason)
				pendingCount--
				if pendingCount == 0 {
					reject(nil, errors)
				}
				return nil, nil
			})
		}
	})
}

// Race 返回一个与第一个完成的 Promise 具有相同状态的 Promise
func Race(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		filtered := make([]*Promise, 0, len(promises))
		for _, p := range promises {
			if p != nil {
				filtered = append(filtered, p)
			}
		}
		promises = filtered

		if len(promises) == 0 {
			resolve(nil, nil)
			return
		}

		isCompleted := false
		var mu sync.Mutex

		for _, promise := range promises {
			promise.Then(func(value interface{}) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				isCompleted = true
				resolve(value, nil)
				return nil, nil
			}, func(reason error) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				if isCompleted {
					return nil, nil
				}
				isCompleted = true
				reject(nil, reason)
				return nil, nil
			})
		}
	})
}
