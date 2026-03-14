package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 1"
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 直接解析为 "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 2"
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 3"
		return value.(string) + " 3", nil
	}, nil)

	// Any() 将等待第一个 promise 被解析，并返回一个带有值的 promise
	result := vl.Any(p1, p2, p3)

	// 从 promise 中获取值并打印
	fmt.Println(">>", result.GetValue().(string))

	// 创建 3 个 promise
	p1 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 被拒绝，原因是 "Promise 1 rejected"
		reject(nil, errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 被拒绝，原因是 "Promise 3 rejected"
		reject(nil, errors.New("Promise 3 rejected"))
	})

	// Any() 将等待所有的 promise 被拒绝，并返回一个带有原因 `AggregateError` 的 promise
	result = vl.Any(p1, p2, p3)

	// 从 promise 中获取原因并打印
	fmt.Println("!!", result.GetReason().Error())
}
