package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 直接解析为 "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，我们将解析的值加上 " vowlink !!"
		return value.(string) + " vowlink !!", nil
	}, func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// 从 promise 中获取值并打印
	fmt.Println("Resolve:", result.GetValue())

	// 这是一个被拒绝的 promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 被拒绝，原因是 "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// 如果 promise 被解析，我们将解析的值加上 " vowlink"
		return value.(string) + " vowlink", nil
	}, func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// 从 promise 中获取拒绝的原因并打印
	fmt.Println("Rejected:", result.GetReason().Error())
}
