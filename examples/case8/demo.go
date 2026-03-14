package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise 1"
		resolve("Promise 1", nil)
	})

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise 3"
		resolve("Promise 3", nil)
	})

	// AllSettled() 将等待所有的 promise 被解析或拒绝，并返回一个带有值的 promise
	result := vl.AllSettled(p1, p2, p3)

	// 从 promise 中获取所有的结果
	for i, r := range result.GetValue().([]interface{}) {
		// 如果结果是一个错误，打印错误信息
		if v, ok := r.(error); ok {
			fmt.Println("!!", i, v.Error())
		} else {
			// 否则，打印结果值
			fmt.Println(">>", i, r.(string))
		}
	}
}
