package vowlink

import "strings"

// PromiseState represents the state of a promise.
// PromiseState 代表一个 Promise 的状态。
type PromiseState uint8

var (
	// defaultOnFulfilledFunc is the default function to be executed when a promise is fulfilled.
	// defaultOnFulfilledFunc 是 Promise 被 fulfilled 时执行的默认函数。
	defaultOnFulfilledFunc = func(value interface{}) (interface{}, error) { return value, nil }

	// defaultOnRejectedFunc is the default function to be executed when a promise is rejected.
	// defaultOnRejectedFunc 是 Promise 被 rejected 时执行的默认函数。
	defaultOnRejectedFunc = func(reason error) (interface{}, error) { return nil, reason }

	// defaultOnFinallyFunc is the default function to be executed when a promise is settled (fulfilled or rejected).
	// defaultOnFinallyFunc 是 Promise 被 settled (fulfilled 或 rejected) 时执行的默认函数。
	defaultOnFinallyFunc = func() {}
)

// AggregateError represents an error that aggregates multiple errors.
// AggregateError 代表一个聚合了多个错误的错误。
type AggregateError struct {
	// Errors 字段是一个 error 类型的切片，用于存储所有的错误
	// Errors field is a slice of error type, used to store all errors
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

// Promise 是一个结构体，它代表一个 Promise。Promise 是一种编程模式，用于处理执行操作。
// Promise is a struct that represents a Promise. A Promise is a programming pattern for handling execution operations.
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

// change 方法用于改变 Promise 的状态，设置 Promise 的值或者拒绝的原因。
// The change method is used to change the state of the Promise, set the value of the Promise or the reason for rejection.
func (p *Promise) change(state PromiseState, value interface{}, reason error) {
	// 只有当 Promise 的状态是 Pending 时，我们才能改变它的状态。
	// We can only change the state of the Promise when its state is Pending.
	if p.state == Pending {
		// 改变 Promise 的状态。
		// Change the state of the Promise.
		p.state = state

		// 设置 Promise 的值，这个值是执行操作的结果。
		// Set the value of the Promise, this value is the result of the execution operation.
		p.value = value

		// 设置 Promise 被拒绝的原因，这个值是执行操作失败的原因。
		// Set the reason the Promise was rejected, this value is the reason the execution operation failed.
		p.reason = reason
	}
}

// resolve 方法用于将 Promise 的状态改变为 Fulfilled，并设置 Promise 的值。
// The resolve method is used to change the state of the Promise to Fulfilled and set the value of the Promise.
func (p *Promise) resolve(value interface{}, reason error) {
	// 调用 change 方法，将 Promise 的状态改变为 Fulfilled，并设置 Promise 的值。
	// Call the change method to change the state of the Promise to Fulfilled and set the value of the Promise.
	p.change(Fulfilled, value, reason)
}

// reject 方法用于将 Promise 的状态改变为 Rejected，并设置 Promise 被拒绝的原因。
// The reject method is used to change the state of the Promise to Rejected and set the reason the Promise was rejected.
func (p *Promise) reject(value interface{}, reason error) {
	// 调用 change 方法，将 Promise 的状态改变为 Rejected，并设置 Promise 被拒绝的原因。
	// Call the change method to change the state of the Promise to Rejected and set the reason the Promise was rejected.
	p.change(Rejected, value, reason)
}

// NewPromise 方法创建一个新的 Promise，并接受一个执行器函数作为参数。执行器函数接受两个参数：resolve 和 reject，分别用于在执行操作成功或失败时改变 Promise 的状态。
// The NewPromise method creates a new Promise and accepts an executor function as a parameter. The executor function accepts two parameters: resolve and reject, which are used to change the state of the Promise when the execution operation succeeds or fails.
func NewPromise(executor func(resolve func(interface{}, error), reject func(interface{}, error))) *Promise {
	// 如果执行器函数是 nil，那么我们没有办法改变 Promise 的状态，所以返回 nil。
	// If the executor function is nil, then we have no way to change the state of the Promise, so return nil.
	if executor == nil {
		return nil // 返回 nil
	}

	// 创建一个新的 Promise，状态为 Pending。Pending 状态表示 Promise 的结果还没有确定。
	// Create a new Promise with the state set to Pending. The Pending state indicates that the result of the Promise is not yet determined.
	p := &Promise{state: Pending}

	// 执行执行器函数，这个函数通常会启动一个执行操作。
	// Execute the executor function, this function usually starts an execution operation.
	executor(p.resolve, p.reject)

	// 返回新创建的 Promise。
	// Return the newly created Promise.
	return p
}

// Then 方法添加一个 fulfilled 和 rejected 处理函数到 Promise。这两个函数分别在 Promise 被 resolve 或 reject 时调用。
// The Then method adds fulfillment and rejection handlers to the Promise. These two functions are called when the Promise is resolved or rejected.
func (p *Promise) Then(onFulfilled func(interface{}) (interface{}, error), onRejected func(error) (interface{}, error)) *Promise {
	// 如果 onFulfilled 函数是 nil，那么我们使用一个默认的 onFulfilled 函数，这个函数直接返回它的参数。
	// If the onFulfilled function is nil, then we use a default onFulfilled function, this function simply returns its parameter.
	if onFulfilled == nil {
		onFulfilled = defaultOnFulfilledFunc
	}

	// 如果 onRejected 函数是 nil 或者 Promise 的拒绝原因是 nil，那么我们使用一个默认的 onRejected 函数，这个函数直接返回它的参数。
	// If the onRejected function is nil or the reason for the Promise's rejection is nil, then we use a default onRejected function, this function simply returns its parameter.
	if onRejected == nil || p.reason == nil {
		onRejected = defaultOnRejectedFunc
	}

	// 返回一个新的 Promise，这个 Promise 的状态和值由 onFulfilled 或 onRejected 函数的返回值决定。
	// Return a new Promise, the state and value of this Promise is determined by the return value of the onFulfilled or onRejected function.
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 根据 Promise 的状态进行不同的操作
		// Perform different operations according to the state of the Promise
		switch p.state {
		// 当 Promise 的状态是 Fulfilled 或 Rejected 时
		// When the state of the Promise is Fulfilled or Rejected
		case Fulfilled, Rejected:
			// 如果 Promise 被拒绝，那么 reason 就不会是 nil
			// If the Promise is rejected, then reason will not be nil
			if p.reason != nil {
				// 调用 onRejected 函数，并将 Promise 被拒绝的原因作为参数传入
				// Call the onRejected function and pass the reason the Promise was rejected as an argument
				reject(onRejected(p.reason))
			} else {
				// 如果 Promise 被解析，那么 value 就是解析的结果
				// If the Promise is resolved, then value is the result of the resolution
				resolve(onFulfilled(p.value))
			}
		}
	})
}

// Catch 方法添加一个 rejected 处理函数到 Promise。这个函数在 Promise 被 reject 时调用。
// The Catch method adds a rejection handler to the Promise. This function is called when the Promise is rejected.
func (p *Promise) Catch(onRejected func(error) (interface{}, error)) *Promise {
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
	return p.Then(func(value interface{}) (interface{}, error) {
		// 先执行 onFinally 函数，无论 Promise 的状态是解析还是拒绝，这个函数都会被调用。
		// First execute the onFinally function, this function will be called regardless of whether the Promise state is resolved or rejected.
		onFinally()

		// 返回原 Promise 的值和被拒绝的原因。如果 Promise 被解析，那么值就是解析的结果，如果 Promise 被拒绝，那么原因就是拒绝的原因。
		// Return the value and the reason for rejection of the original Promise. If the Promise is resolved, the value is the result of the resolution. If the Promise is rejected, the reason is the reason for the rejection.
		return value, p.reason
	}, func(reason error) (interface{}, error) {
		// 先执行 onFinally 函数，无论 Promise 的状态是解析还是拒绝，这个函数都会被调用。
		// First execute the onFinally function, this function will be called regardless of whether the Promise state is resolved or rejected.
		onFinally()

		// 返回原 Promise 的值和被拒绝的原因。如果 Promise 被解析，那么值就是解析的结果，如果 Promise 被拒绝，那么原因就是拒绝的原因。
		// Return the value and the reason for rejection of the original Promise. If the Promise is resolved, the value is the result of the resolution. If the Promise is rejected, the reason is the reason for the rejection.
		return p.value, reason
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
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
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
			promise.Then(func(value interface{}) (interface{}, error) {
				// 存储 Promise 的结果到对应的位置。
				// Store the result of the Promise at the corresponding position.
				values[i] = value

				// 增加已完成的 Promise 的数量。
				// Increase the number of fulfilled Promises.
				count++

				// 如果所有的 Promise 都已完成，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are fulfilled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values, nil)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			}, func(reason error) (interface{}, error) {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise，原因为被拒绝的 Promise 的原因。
				// If any Promise is rejected, reject this Promise with the reason of the rejected Promise.
				reject(nil, reason)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			})
		}
	})
}

// AllSettled 方法返回一个新的 Promise，它在所有给定的 Promise 都 settled (fulfilled 或 rejected) 时解析。
// The AllSettled method returns a new Promise that resolves when all the given Promises are settled (fulfilled or rejected).
func AllSettled(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
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
			promise.Then(func(value interface{}) (interface{}, error) {
				// 存储 Promise 的结果到对应的位置。
				// Store the result of the Promise at the corresponding position.
				values[i] = value

				// 增加已 settled 的 Promise 的数量。
				// Increase the number of settled Promises.
				count++

				// 如果所有的 Promise 都已 settled，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are settled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values, nil)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			}, func(reason error) (interface{}, error) {
				// 存储 Promise 被拒绝的原因到对应的位置。
				// Store the reason the Promise was rejected at the corresponding position.
				values[i] = reason

				// 增加已 settled 的 Promise 的数量。
				// Increase the number of settled Promises.
				count++

				// 如果所有的 Promise 都已 settled，那么解析这个 Promise，值为所有 Promise 的结果。
				// If all Promises are settled, resolve this Promise with the results of all Promises.
				if count == len(promises) {
					resolve(values, nil)
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			})
		}
	})
}

// Any 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成时解析，或者在所有的 Promise 都被拒绝时拒绝。
// The Any method returns a new Promise that resolves when any of the given Promises is fulfilled, or rejects when all of the Promises are rejected.
func Any(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
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
			promise.Then(func(value interface{}) (interface{}, error) {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise，值为完成的 Promise 的值。
				// If any Promise is fulfilled, resolve this Promise with the value of the fulfilled Promise.
				resolve(value, nil)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			}, func(reason error) (interface{}, error) {
				// 存储 Promise 被拒绝的原因到对应的位置。
				// Store the reason the Promise was rejected at the corresponding position.
				errors[i] = reason

				// 增加已拒绝的 Promise 的数量。
				// Increase the number of rejected Promises.
				count++

				// 如果所有的 Promise 都被拒绝，那么拒绝这个 Promise，原因为所有 Promise 被拒绝的原因的集合。
				// If all Promises are rejected, reject this Promise with the collection of reasons all Promises were rejected.
				if count == len(promises) {
					reject(nil, &AggregateError{Errors: errors})
				}

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			})
		}
	})
}

// Race 方法返回一个新的 Promise，它在任何一个给定的 Promise 完成或被拒绝时解析或拒绝。
// The Race method returns a new Promise that resolves or rejects as soon as any of the given Promises is fulfilled or rejected.
func Race(promises ...*Promise) *Promise {
	// 创建一个新的 Promise，它的解析和拒绝函数由下面的逻辑决定。
	// Create a new Promise, its resolve and reject functions are determined by the following logic.
	return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 遍历给定的每一个 Promise。
		// Iterate over each given Promise.
		for _, promise := range promises {
			// 对每个 Promise 添加一个成功和失败的回调。
			// Add a success and failure callback to each Promise.
			promise.Then(func(value interface{}) (interface{}, error) {
				// 如果有任何一个 Promise 完成，那么解析这个 Promise，值为完成的 Promise 的值。
				// If any Promise is fulfilled, resolve this Promise with the value of the fulfilled Promise.
				resolve(value, nil)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			}, func(reason error) (interface{}, error) {
				// 如果有任何一个 Promise 被拒绝，那么拒绝这个 Promise，原因为被拒绝的 Promise 的原因。
				// If any Promise is rejected, reject this Promise with the reason of the rejected Promise.
				reject(nil, reason)

				// 返回 nil，因为我们不需要在这里改变 Promise 的值。
				// Return nil because we don't need to change the value of the Promise here.
				return nil, nil
			})
		}
	})
}
