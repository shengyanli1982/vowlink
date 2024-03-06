package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3"
	}, nil)

	// Any() 将等待第一个 promise 被解析，并返回一个带有值的 promise
	// Any() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Any(p1, p2, p3)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(">>", result.GetValue().(string))

	// 创建 3 个 promise
	// Create 3 promises
	p1 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 被拒绝，原因是 "Promise 1 rejected"
		// The first promise is rejected with the reason "Promise 1 rejected"
		reject(errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 被拒绝，原因是 "Promise 3 rejected"
		// The third promise is rejected with the reason "Promise 3 rejected"
		reject(errors.New("Promise 3 rejected"))
	})

	// Any() 将等待所有的 promise 被拒绝，并返回一个带有原因 `AggregateError` 的 promise
	// Any() will wait for all promises to be rejected, and return a promise with the reason `AggregateError`
	result = vl.Any(p1, p2, p3)

	// 从 promise 中获取原因并打印
	// Get the reason from the promise and print it
	fmt.Println("!!", result.GetReason().Error())
}
