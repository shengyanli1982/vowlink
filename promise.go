package vowlink

import "strings"

// PromiseState represents the state of a promise.
// PromiseState 代表一个 Promise 的状态。
type PromiseState uint8

var (
	// defaultOnFulfilledFunc is the default function to be executed when a promise is fulfilled.
	// defaultOnFulfilledFunc 是 Promise 被 fulfilled 时执行的默认函数。
	defaultOnFulfilledFunc = func(value interface{}) interface{} { return value }

	// defaultOnRejectedFunc is the default function to be executed when a promise is rejected.
	// defaultOnRejectedFunc 是 Promise 被 rejected 时执行的默认函数。
	defaultOnRejectedFunc = func(reason error) error { return reason }

	// defaultOnFinallyFunc is the default function to be executed when a promise is settled (fulfilled or rejected).
	// defaultOnFinallyFunc 是 Promise 被 settled (fulfilled 或 rejected) 时执行的默认函数。
	defaultOnFinallyFunc = func() {}
)

// AggregateError represents an error that aggregates multiple errors.
// AggregateError 代表一个聚合了多个错误的错误。
type AggregateError struct {
	Errors []error
}

// Error 方法返回 AggregateError 的字符串表示。
// The Error method returns the string representation of the AggregateError.
func (ae *AggregateError) Error() string {
	// 创建一个切片来存储所有错误的字符串表示
	// Create a slice to store the string representations of all errors
	errStrings := make([]string, len(ae.Errors))
	for i, err := range ae.Errors {
		// 将每个错误转换为字符串并存储在切片中
		// Convert each error to a string and store it in the slice
		errStrings[i] = err.Error()
	}
	// 返回所有错误的字符串表示，用逗号分隔
	// Return the string representations of all errors, separated by commas
	return "All promises were rejected: " + strings.Join(errStrings, ", ")
}

// 定义 Promise 的三种状态：Pending、Fulfilled 和 Rejected
// Define the three states of a Promise: Pending, Fulfilled, and Rejected
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
// Promise is a struct that represents a Promise.
type Promise struct {
	state  PromiseState // Promise 的状态
	value  interface{}  // Promise 的值
	reason error        // Promise 被拒绝的原因
}

// resolve 方法将 Promise 的状态设置为 Fulfilled 并设置值。
// The resolve method sets the state of the Promise to Fulfilled and sets the value.
func (p *Promise) resolve(value interface{}) {
	if p.state == Pending { // 如果 Promise 的状态是 Pending
		p.state = Fulfilled // 将状态设置为 Fulfilled
		p.value = value     // 设置 Promise 的值
	}
}

// reject 方法将 Promise 的状态设置为 Rejected 并设置原因。
// The reject method sets the state of the Promise to Rejected and sets the reason.
func (p *Promise) reject(reason error) {
	if p.state == Pending { // 如果 Promise 的状态是 Pending
		p.state = Rejected // 将状态设置为 Rejected
		p.reason = reason  // 设置 Promise 被拒绝的原因
	}
}

// NewPromise 方法创建一个新的 Promise，并接受一个执行器函数作为参数。
// The NewPromise method creates a new Promise and accepts an executor function as a parameter.
func NewPromise(executor func(resolve func(interface{}), reject func(error))) *Promise {
	// 如果执行器函数是 nil
	if executor == nil {
		return nil // 返回 nil
	}

	// 创建一个新的 Promise，状态为 Pending
	p := &Promise{state: Pending}

	// 执行执行器函数
	executor(p.resolve, p.reject)

	// 返回新创建的 Promise
	return p
}

// Then 方法添加一个 fulfilled 和 rejected 处理函数到 Promise。
// The Then method adds fulfillment and rejection handlers to the Promise.
func (p *Promise) Then(onFulfilled func(interface{}) interface{}, onRejected func(error) error) *Promise {
	if onFulfilled == nil { // 如果 onFulfilled 函数是 nil
		onFulfilled = defaultOnFulfilledFunc // 使用默认的 onFulfilled 函数
	}
	if onRejected == nil { // 如果 onRejected 函数是 nil
		onRejected = defaultOnRejectedFunc // 使用默认的 onRejected 函数
	}

	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		switch p.state { // 根据 Promise 的状态
		case Fulfilled: // 如果 Promise 的状态是 Fulfilled
			resolve(onFulfilled(p.value)) // 执行 onFulfilled 函数并解析 Promise
		case Rejected: // 如果 Promise 的状态是 Rejected
			reject(onRejected(p.reason)) // 执行 onRejected 函数并拒绝 Promise
		}
	})
}

// Catch 添加一个 rejected 处理函数到 Promise。
// Catch adds a rejection handler to the promise.
func (p *Promise) Catch(onRejected func(error) error) *Promise {
	// 将 onRejected 函数添加到 Promise，如果 Promise 被 reject，则执行这个函数。
	// Add the onRejected function to the Promise, if the Promise is rejected, this function will be executed.
	return p.Then(nil, onRejected)
}

// Finally 添加一个 settled (fulfilled 或 rejected) 处理函数到 Promise。
// Finally adds a finally handler to the promise.
func (p *Promise) Finally(onFinally func()) *Promise {
	// 如果 onFinally 是 nil，则将其替换为默认函数。
	// If onFinally is nil, replace it with the default function.
	if onFinally == nil {
		onFinally = defaultOnFinallyFunc
	}

	// 返回一个新的 Promise，该 Promise 在 p settled 时执行 onFinally。
	// Return a new promise that executes onFinally when p is settled.
	return p.Then(func(value interface{}) interface{} {
		onFinally()
		return value
	}, func(reason error) error {
		onFinally()
		return reason
	})
}

// GetValue 方法返回 Promise 的值。
// The GetValue method returns the value of the Promise.
func (p *Promise) GetValue() interface{} {
	// 返回 Promise 的值
	// Return the value of the Promise
	return p.value
}

// GetReason 方法返回 Promise 的原因。
// The GetReason method returns the reason of the Promise.
func (p *Promise) GetReason() error {
	// 返回 Promise 的原因
	// Return the reason of the Promise
	return p.reason
}

// All 方法返回一个新的 Promise，它在所有给定的 Promise 都完成时解析，或者在任何一个 Promise 被拒绝时拒绝。
// The All method returns a new Promise that resolves when all the given Promises are fulfilled, or rejects when any of the Promises is rejected.
func All(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的结果
		// Create a slice to store the results of all Promises
		values := make([]interface{}, len(promises))
		// 创建一个计数器来跟踪已完成的 Promise 的数量
		// Create a counter to track the number of fulfilled Promises
		count := 0
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调
			// Add a success and failure callback to each Promise
			promise.Then(func(value interface{}) interface{} {
				// 存储 Promise 的结果
				// Store the result of the Promise
				values[i] = value
				count++
				// 如果所有的 Promise 都已完成，那么解析这个 Promise
				// If all Promises are fulfilled, resolve this Promise
				if count == len(promises) {
					resolve(values)
				}
				return nil
			}, func(reason error) error {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise
				// If any Promise is rejected, reject this Promise
				reject(reason)
				return nil
			})
		}
	})
}

// AllSettled 方法返回一个新的 Promise，它在所有给定的 Promise 都 settled (fulfilled 或 rejected) 时解析。
// The AllSettled method returns a new Promise that resolves when all the given Promises are settled (fulfilled or rejected).
func AllSettled(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的结果
		// Create a slice to store the results of all Promises
		values := make([]interface{}, len(promises))
		// 创建一个计数器来跟踪已 settled 的 Promise 的数量
		// Create a counter to track the number of settled Promises
		count := 0
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调
			// Add a success and failure callback to each Promise
			promise.Then(func(value interface{}) interface{} {
				// 存储 Promise 的结果
				// Store the result of the Promise
				values[i] = value
				count++
				// 如果所有的 Promise 都已 settled，那么解析这个 Promise
				// If all Promises are settled, resolve this Promise
				if count == len(promises) {
					resolve(values)
				}
				return nil
			}, func(reason error) error {
				// 存储 Promise 的原因
				// Store the reason of the Promise
				values[i] = reason
				count++
				// 如果所有的 Promise 都已 settled，那么解析这个 Promise
				// If all Promises are settled, resolve this Promise
				if count == len(promises) {
					resolve(values)
				}
				return nil
			})
		}
	})
}

// Any 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成时解析，或者在所有的 Promise 都被拒绝时拒绝。
// The Any method returns a new Promise that resolves when any of the given Promises is fulfilled, or rejects when all of the Promises are rejected.
func Any(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的错误
		// Create a slice to store the errors of all Promises
		errors := make([]error, len(promises))
		// 创建一个计数器来跟踪已拒绝的 Promise 的数量
		// Create a counter to track the number of rejected Promises
		count := 0
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调
			// Add a success and failure callback to each Promise
			promise.Then(func(value interface{}) interface{} {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise
				// If any Promise is fulfilled, resolve this Promise
				resolve(value)
				return nil
			}, func(reason error) error {
				// 存储 Promise 的错误
				// Store the error of the Promise
				errors[i] = reason
				count++
				// 如果所有的 Promise 都被拒绝，那么拒绝这个 Promise
				// If all Promises are rejected, reject this Promise
				if count == len(promises) {
					reject(&AggregateError{Errors: errors})
				}
				return nil
			})
		}
	})
}

// Race 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成或被拒绝时解析或拒绝。
// The Race method returns a new Promise that resolves or rejects as soon as any of the given Promises is fulfilled or rejected.
func Race(promises ...*Promise) *Promise {
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		for _, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调
			// Add a success and failure callback to each Promise
			promise.Then(func(value interface{}) interface{} {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise
				// If any Promise is fulfilled, resolve this Promise
				resolve(value)
				return nil
			}, func(reason error) error {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise
				// If any Promise is rejected, reject this Promise
				reject(reason)
				return nil
			})
		}
	})
}
