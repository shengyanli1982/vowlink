package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// All() 将等待所有的 promise 被解析，并返回一个带有所有值的 promise
	// All() will wait for all promises to be resolved, and return a promise with all the values
	result := vl.All(p1, p2, p3)

	// 从 promise 中获取所有的值并打印
	// Get all the values from the promise and print them
	for i, str := range result.GetValue().([]interface{}) {
		fmt.Println(">>", i, str.(string))
	}
}
