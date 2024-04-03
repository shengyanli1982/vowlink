package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise 1"
		// The first promise is directly resolved to "Promise 1"
		resolve("Promise 1", nil)
	})

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise 3"
		// The third promise is directly resolved to "Promise 3"
		resolve("Promise 3", nil)
	})

	// AllSettled() 将等待所有的 promise 被解析或拒绝，并返回一个带有值的 promise
	// AllSettled() will wait for all promises to be resolved or rejected, and return a promise with the value
	result := vl.AllSettled(p1, p2, p3)

	// 从 promise 中获取所有的结果
	// Get all the results from the promise
	for i, r := range result.GetValue().([]interface{}) {
		// 如果结果是一个错误，打印错误信息
		// If the result is an error, print the error message
		if v, ok := r.(error); ok {
			fmt.Println("!!", i, v.Error())
		} else {
			// 否则，打印结果值
			// Otherwise, print the result value
			fmt.Println(">>", i, r.(string))
		}
	}
}
