package vowlink

import "strings"

// PromiseState 代表一个 Promise 的状态。
// PromiseState represents the state of a promise.
type PromiseState uint8

var (
	// defaultOnFulfilledFunc 是 Promise 被 fulfilled 时执行的默认函数。
	// defaultOnFulfilledFunc is the default function to be executed when a promise is fulfilled.
	defaultOnFulfilledFunc = func(value interface{}) interface{} { return value }

	// defaultOnRejectedFunc 是 Promise 被 rejected 时执行的默认函数。
	// defaultOnRejectedFunc is the default function to be executed when a promise is rejected.
	defaultOnRejectedFunc = func(reason error) error { return reason }

	// defaultOnFinallyFunc 是 Promise 被 settled (fulfilled 或 rejected) 时执行的默认函数。
	// defaultOnFinallyFunc is the default function to be executed when a promise is settled (fulfilled or rejected).
	defaultOnFinallyFunc = func() {}
)

// AggregateError 代表一个聚合了多个错误的错误。
// AggregateError represents an error that aggregates multiple errors.
type AggregateError struct {
	Errors []error
}

// Error 返回 AggregateError 的字符串表示。
// Error returns the string representation of the AggregateError.
func (ae *AggregateError) Error() string {
	errStrings := make([]string, len(ae.Errors))
	for i, err := range ae.Errors {
		errStrings[i] = err.Error()
	}
	return "All promises were rejected: " + strings.Join(errStrings, ", ")
}

const (
	// Pending 代表一个 Promise 的 pending 状态。
	// Pending represents the pending state of a promise.
	Pending PromiseState = iota

	// Fulfilled 代表一个 Promise 的 fulfilled 状态。
	// Fulfilled represents the fulfilled state of a promise.
	Fulfilled

	// Rejected 代表一个 Promise 的 rejected 状态。
	// Rejected represents the rejected state of a promise.
	Rejected
)

// Promise 是一个结构体，代表一个 Promise。
// Promise represents a promise.
type Promise struct {
	state  PromiseState
	value  interface{}
	reason error
}

// resolve 设置 Promise 的状态为 Fulfilled 并设置值。
// resolve sets the promise state to Fulfilled and sets the value.
func (p *Promise) resolve(value interface{}) {
	if p.state == Pending {
		p.state = Fulfilled
		p.value = value
	}
}

// reject 设置 Promise 的状态为 Rejected 并设置原因。
// reject sets the promise state to Rejected and sets the reason.
func (p *Promise) reject(reason error) {
	if p.state == Pending {
		p.state = Rejected
		p.reason = reason
	}
}

// NewPromise 创建一个新的 Promise。
// NewPromise creates a new promise with the given executor function.
func NewPromise(executor func(resolve func(interface{}), reject func(error))) *Promise {
	if executor == nil {
		return nil
	}
	p := &Promise{state: Pending}
	executor(p.resolve, p.reject)
	return p
}

// Then 添加一个 fulfilled 和 rejected 处理函数到 Promise。
// Then adds fulfillment and rejection handlers to the promise.
func (p *Promise) Then(onFulfilled func(interface{}) interface{}, onRejected func(error) error) *Promise {
	if onFulfilled == nil {
		onFulfilled = defaultOnFulfilledFunc
	}
	if onRejected == nil {
		onRejected = defaultOnRejectedFunc
	}

	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		switch p.state {
		case Fulfilled:
			resolve(onFulfilled(p.value))
		case Rejected:
			reject(onRejected(p.reason))
		}
	})
}

// Catch 添加一个 rejected 处理函数到 Promise。
// Catch adds a rejection handler to the promise.
func (p *Promise) Catch(onRejected func(error) error) *Promise {
	return p.Then(nil, onRejected)
}

// Finally 添加一个 settled (fulfilled 或 rejected) 处理函数到 Promise。
// Finally adds a finally handler to the promise.
func (p *Promise) Finally(onFinally func()) *Promise {
	if onFinally == nil {
		onFinally = defaultOnFinallyFunc
	}
	return p.Then(func(value interface{}) interface{} {
		onFinally()
		return value
	}, func(reason error) error {
		onFinally()
		return reason
	})
}

// GetValue 返回 Promise 的值。
// GetValue returns the value of the promise.
func (p *Promise) GetValue() interface{} {
	return p.value
}

// GetReason 返回 Promise 的原因。
// GetReason returns the reason of the promise.
func (p *Promise) GetReason() error {
	return p.reason
}

// All 返回一个 Promise，当所有输入的 Promise 都 fulfilled 时，该 Promise fulfilled，或者当任何一个 Promise rejected 时，该 Promise rejected。
// All returns a promise that resolves when all the input promises are fulfilled, or rejects if any of the input promises are rejected.
func All(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		values := make([]interface{}, len(promises))
		count := 0
		for i, promise := range promises {
			promise.Then(func(value interface{}) interface{} {
				values[i] = value
				count++
				if count == len(promises) {
					resolve(values)
				}
				return nil
			}, func(reason error) error {
				reject(reason)
				return nil
			})
		}
	})
}

// AllSettled 返回一个 Promise，当所有输入的 Promise 都 settled (fulfilled 或 rejected) 时，该 Promise fulfilled。
// AllSettled returns a promise that resolves when all the input promises are settled (fulfilled or rejected).
func AllSettled(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		values := make([]interface{}, len(promises))
		count := 0
		for i, promise := range promises {
			promise.Then(func(value interface{}) interface{} {
				values[i] = value
				count++
				if count == len(promises) {
					resolve(values)
				}
				return nil
			}, func(reason error) error {
				values[i] = reason
				count++
				if count == len(promises) {
					resolve(values)
				}
				return nil
			})
		}
	})
}

// Any 返回一个 Promise，当任何一个输入的 Promise fulfilled 时，该 Promise fulfilled，或者当所有的 Promise 都 rejected 时，该 Promise rejected。
// Any returns a promise that resolves when any of the input promises is fulfilled, or rejects if all of the input promises are rejected.
func Any(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		errors := make([]error, len(promises))
		count := 0
		for i, promise := range promises {
			promise.Then(func(value interface{}) interface{} {
				resolve(value)
				return nil
			}, func(reason error) error {
				errors[i] = reason
				count++
				if count == len(promises) {
					reject(&AggregateError{Errors: errors})
				}
				return nil
			})
		}
	})
}

// Race 返回一个 Promise，当任何一个输入的 Promise fulfilled 或 rejected 时，该 Promise settled。
// Race returns a promise that resolves or rejects as soon as any of the input promises resolves or rejects.
func Race(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		for _, promise := range promises {
			promise.Then(func(value interface{}) interface{} {
				resolve(value)
				return nil
			}, func(reason error) error {
				reject(reason)
				return nil
			})
		}
	})
}
