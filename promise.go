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
	// 创建一个切片来存储所有错误的字符串表示。这个切片的长度等于错误的数量。
	// Create a slice to store the string representations of all errors. The length of this slice is equal to the number of errors.
	errStrings := make([]string, len(ae.Errors))

	// 遍历所有的错误。
	// Iterate over all errors.
	for i, err := range ae.Errors {
		// 将每个错误转换为字符串并存储在切片中。
		// Convert each error to a string and store it in the slice.
		errStrings[i] = err.Error()
	}

	// 返回所有错误的字符串表示，用逗号分隔。这个字符串表示了所有 Promise 被拒绝的原因。
	// Return the string representations of all errors, separated by commas. This string represents the reasons all Promises were rejected.
	return "All promises were rejected: " + strings.Join(errStrings, ", ")
}

// 定义 Promise 的三种状态：Pending、Fulfilled 和 Rejected
// Define the three states of a Promise: Pending, Fulfilled, and Rejected
const (
	// Pending 代表一个 Promise 的 pending 状态，即 Promise 还在等待结果，既没有完成也没有被拒绝。
	// Pending represents the pending state of a promise, which means the Promise is still waiting for the result, neither fulfilled nor rejected.
	Pending PromiseState = iota

	// Fulfilled 代表一个 Promise 的 fulfilled 状态，即 Promise 已经完成，有一个确定的结果值。
	// Fulfilled represents the fulfilled state of a promise, which means the Promise is fulfilled and has a definite result value.
	Fulfilled

	// Rejected 代表一个 Promise 的 rejected 状态，即 Promise 已经被拒绝，有一个确定的错误原因。
	// Rejected represents the rejected state of a promise, which means the Promise is rejected and has a definite reason for the error.
	Rejected
)

// Promise 是一个结构体，它代表一个 Promise。Promise 是一种编程模式，用于处理异步操作。
// Promise is a struct that represents a Promise. A Promise is a programming pattern for handling asynchronous operations.
type Promise struct {
	// state 是 Promise 的状态，它可以是 Pending、Fulfilled 或 Rejected。
	// state is the state of the Promise, it can be Pending, Fulfilled, or Rejected.
	state PromiseState

	// value 是 Promise 的值，当 Promise 被 resolve 时，这个值会被设置。
	// value is the value of the Promise, this value will be set when the Promise is resolved.
	value interface{}

	// reason 是 Promise 被拒绝的原因，当 Promise 被 reject 时，这个值会被设置。
	// reason is the reason the Promise was rejected, this value will be set when the Promise is rejected.
	reason error
}

// resolve 是一个方法，它将 Promise 的状态设置为 Fulfilled 并设置值。这个方法通常在异步操作成功完成时调用。
// resolve is a method that sets the state of the Promise to Fulfilled and sets the value. This method is typically called when the asynchronous operation successfully completes.
func (p *Promise) resolve(value interface{}) {
	// 如果 Promise 的状态是 Pending，那么我们可以将其状态设置为 Fulfilled。
	// If the state of the Promise is Pending, then we can set its state to Fulfilled.
	if p.state == Pending {
		// 将状态设置为 Fulfilled。
		// Set the state to Fulfilled.
		p.state = Fulfilled

		// 设置 Promise 的值，这个值是异步操作的结果。
		// Set the value of the Promise, this value is the result of the asynchronous operation.
		p.value = value
	}
}

// reject 是一个方法，它将 Promise 的状态设置为 Rejected 并设置原因。这个方法通常在异步操作失败时调用。
// reject is a method that sets the state of the Promise to Rejected and sets the reason. This method is typically called when the asynchronous operation fails.
func (p *Promise) reject(reason error) {
	// 如果 Promise 的状态是 Pending，那么我们可以将其状态设置为 Rejected。
	// If the state of the Promise is Pending, then we can set its state to Rejected.
	if p.state == Pending {
		// 将状态设置为 Rejected。
		// Set the state to Rejected.
		p.state = Rejected

		// 设置 Promise 被拒绝的原因，这个值是异步操作失败的原因。
		// Set the reason the Promise was rejected, this value is the reason the asynchronous operation failed.
		p.reason = reason
	}
}

// NewPromise 方法创建一个新的 Promise，并接受一个执行器函数作为参数。执行器函数接受两个参数：resolve 和 reject，分别用于在异步操作成功或失败时改变 Promise 的状态。
// The NewPromise method creates a new Promise and accepts an executor function as a parameter. The executor function accepts two parameters: resolve and reject, which are used to change the state of the Promise when the asynchronous operation succeeds or fails.
func NewPromise(executor func(resolve func(interface{}), reject func(error))) *Promise {
	// 如果执行器函数是 nil，那么我们没有办法改变 Promise 的状态，所以返回 nil。
	// If the executor function is nil, then we have no way to change the state of the Promise, so return nil.
	if executor == nil {
		return nil // 返回 nil
	}

	// 创建一个新的 Promise，状态为 Pending。Pending 状态表示 Promise 的结果还没有确定。
	// Create a new Promise with the state set to Pending. The Pending state indicates that the result of the Promise is not yet determined.
	p := &Promise{state: Pending}

	// 执行执行器函数，这个函数通常会启动一个异步操作。
	// Execute the executor function, this function usually starts an asynchronous operation.
	executor(p.resolve, p.reject)

	// 返回新创建的 Promise。
	// Return the newly created Promise.
	return p
}

// Then 方法添加一个 fulfilled 和 rejected 处理函数到 Promise。这两个函数分别在 Promise 被 resolve 或 reject 时调用。
// The Then method adds fulfillment and rejection handlers to the Promise. These two functions are called when the Promise is resolved or rejected.
func (p *Promise) Then(onFulfilled func(interface{}) interface{}, onRejected func(error) error) *Promise {
	// 如果 onFulfilled 函数是 nil，那么我们使用一个默认的 onFulfilled 函数，这个函数直接返回它的参数。
	// If the onFulfilled function is nil, then we use a default onFulfilled function, this function simply returns its parameter.
	if onFulfilled == nil {
		onFulfilled = defaultOnFulfilledFunc
	}
	// 如果 onRejected 函数是 nil，那么我们使用一个默认的 onRejected 函数，这个函数直接返回它的参数。
	// If the onRejected function is nil, then we use a default onRejected function, this function simply returns its parameter.
	if onRejected == nil {
		onRejected = defaultOnRejectedFunc
	}

	// 返回一个新的 Promise，这个 Promise 的状态和值由 onFulfilled 或 onRejected 函数的返回值决定。
	// Return a new Promise, the state and value of this Promise is determined by the return value of the onFulfilled or onRejected function.
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 根据 Promise 的状态
		// According to the state of the Promise
		switch p.state {

		// 如果 Promise 的状态是 Fulfilled
		// If the state of the Promise is Fulfilled
		case Fulfilled:
			// 执行 onFulfilled 函数并解析 Promise。onFulfilled 函数的参数是原 Promise 的值。
			// Execute the onFulfilled function and resolve the Promise. The parameter of the onFulfilled function is the value of the original Promise.
			resolve(onFulfilled(p.value))

		// 如果 Promise 的状态是 Rejected
		// If the state of the Promise is Rejected
		case Rejected:
			// 执行 onRejected 函数并拒绝 Promise。onRejected 函数的参数是原 Promise 被拒绝的原因。
			// Execute the onRejected function and reject the Promise. The parameter of the onRejected function is the reason the original Promise was rejected.
			reject(onRejected(p.reason))
		}
	})
}

// Catch 方法添加一个 rejected 处理函数到 Promise。这个函数在 Promise 被 reject 时调用。
// The Catch method adds a rejection handler to the Promise. This function is called when the Promise is rejected.
func (p *Promise) Catch(onRejected func(error) error) *Promise {
	// 将 onRejected 函数添加到 Promise，如果 Promise 被 reject，则执行这个函数。
	// Add the onRejected function to the Promise, if the Promise is rejected, this function will be executed.
	return p.Then(nil, onRejected)
}

// Finally 方法添加一个 settled (fulfilled 或 rejected) 处理函数到 Promise。这个函数在 Promise 被 resolve 或 reject 时调用。
// The Finally method adds a finally handler to the Promise. This function is called when the Promise is resolved or rejected.
func (p *Promise) Finally(onFinally func()) *Promise {
	// 如果 onFinally 是 nil，则将其替换为默认函数。这个默认函数什么也不做。
	// If onFinally is nil, replace it with the default function. This default function does nothing.
	if onFinally == nil {
		onFinally = defaultOnFinallyFunc
	}

	// 返回一个新的 Promise，该 Promise 在 p settled 时执行 onFinally。无论 Promise 被 resolve 还是 reject，onFinally 都会被调用。
	// Return a new Promise that executes onFinally when p is settled. Whether the Promise is resolved or rejected, onFinally will be called.
	return p.Then(func(value interface{}) interface{} {
		// 先执行 onFinally 函数，然后返回原 Promise 的值。
		// First execute the onFinally function, then return the value of the original Promise.
		onFinally()

		// 返回原 Promise 的值。这个值将被传递给下一个 Then 或 Catch 方法。
		// Return the value of the original Promise. This value will be passed to the next Then or Catch method.
		return value
	}, func(reason error) error {
		// 先执行 onFinally 函数，然后返回原 Promise 被拒绝的原因。
		// First execute the onFinally function, then return the reason the original Promise was rejected.
		onFinally()

		// 先执行 onFinally 函数，然后返回原 Promise 的值。
		// First execute the onFinally function, then return the value of the original Promise.
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
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的结果。这个切片的长度等于给定的 Promise 的数量。
		// Create a slice to store the results of all Promises. The length of this slice is equal to the number of given Promises.
		values := make([]interface{}, len(promises))

		// 创建一个计数器来跟踪已完成的 Promise 的数量。初始值为 0。
		// Create a counter to track the number of fulfilled Promises. The initial value is 0.
		count := 0

		// 遍历给定的每一个 Promise。
		// Iterate over each given Promise.
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调。
			// Add a success and failure callback to each Promise.
			promise.Then(func(value interface{}) interface{} {
				// 存储 Promise 的结果到对应的位置。
				// Store the result of the Promise at the corresponding position.
				values[i] = value

				// 增加已完成的 Promise 的数量。
				// Increase the number of fulfilled Promises.
				count++

				// 如果所有的 Promise 都已完成，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are fulfilled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			}, func(reason error) error {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise，原因为被拒绝的 Promise 的原因。
				// If any Promise is rejected, reject this Promise with the reason of the rejected Promise.
				reject(reason)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			})
		}
	})
}

// AllSettled 方法返回一个新的 Promise，它在所有给定的 Promise 都 settled (fulfilled 或 rejected) 时解析。
// The AllSettled method returns a new Promise that resolves when all the given Promises are settled (fulfilled or rejected).
func AllSettled(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的结果。这个切片的长度等于给定的 Promise 的数量。
		// Create a slice to store the results of all Promises. The length of this slice is equal to the number of given Promises.
		values := make([]interface{}, len(promises))

		// 创建一个计数器来跟踪已 settled 的 Promise 的数量。初始值为 0。
		// Create a counter to track the number of settled Promises. The initial value is 0.
		count := 0

		// 遍历给定的每一个 Promise。
		// Iterate over each given Promise.
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调。
			// Add a success and failure callback to each Promise.
			promise.Then(func(value interface{}) interface{} {
				// 存储 Promise 的结果到对应的位置。
				// Store the result of the Promise at the corresponding position.
				values[i] = value

				// 增加已 settled 的 Promise 的数量。
				// Increase the number of settled Promises.
				count++

				// 如果所有的 Promise 都已 settled，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are settled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			}, func(reason error) error {
				// 存储 Promise 被拒绝的原因到对应的位置。
				// Store the reason the Promise was rejected at the corresponding position.
				values[i] = reason

				// 增加已 settled 的 Promise 的数量。
				// Increase the number of settled Promises.
				count++

				// 如果所有的 Promise 都已 settled，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are settled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			})
		}
	})
}

// Any 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成时解析，或者在所有的 Promise 都被拒绝时拒绝。
// The Any method returns a new Promise that resolves when any of the given Promises is fulfilled, or rejects when all of the Promises are rejected.
func Any(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 创建一个切片来存储所有 Promise 的错误。这个切片的长度等于给定的 Promise 的数量。
		// Create a slice to store the errors of all Promises. The length of this slice is equal to the number of given Promises.
		errors := make([]error, len(promises))

		// 创建一个计数器来跟踪已拒绝的 Promise 的数量。初始值为 0。
		// Create a counter to track the number of rejected Promises. The initial value is 0.
		count := 0

		// 遍历给定的每一个 Promise。
		// Iterate over each given Promise.
		for i, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调。
			// Add a success and failure callback to each Promise.
			promise.Then(func(value interface{}) interface{} {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise，值为完成的 Promise 的值。
				// If any Promise is fulfilled, resolve this Promise with the value of the fulfilled Promise.
				resolve(value)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			}, func(reason error) error {
				// 存储 Promise 被拒绝的原因到对应的位置。
				// Store the reason the Promise was rejected at the corresponding position.
				errors[i] = reason

				// 增加已拒绝的 Promise 的数量。
				// Increase the number of rejected Promises.
				count++

				// 如果所有的 Promise 都被拒绝，那么拒绝这个 Promise，原因为所有 Promise 被拒绝的原因的集合。
				// If all Promises are rejected, reject this Promise with the collection of reasons all Promises were rejected.
				if count == len(promises) {
					reject(&AggregateError{Errors: errors})
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			})
		}
	})
}

// Race 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成或被拒绝时解析或拒绝。
// The Race method returns a new Promise that resolves or rejects as soon as any of the given Promises is fulfilled or rejected.
func Race(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 遍历给定的每一个 Promise。
		// Iterate over each given Promise.
		for _, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调。
			// Add a success and failure callback to each Promise.
			promise.Then(func(value interface{}) interface{} {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise，值为完成的 Promise 的值。
				// If any Promise is fulfilled, resolve this Promise with the value of the fulfilled Promise.
				resolve(value)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			}, func(reason error) error {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise，原因为被拒绝的 Promise 的原因。
				// If any Promise is rejected, reject this Promise with the reason of the rejected Promise.
				reject(reason)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil
			})
		}
	})
}
