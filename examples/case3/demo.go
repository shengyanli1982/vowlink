package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情。
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，我们将解析的值加上 " Copilot !!"
		// In the Then method, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	}).Finally(func() {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 1"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 1"
		fmt.Println("finally 1")
	})

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		// 如果 promise 被解析，我们将解析的值加上 " Copilot !!"
		// If the promise is resolved, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	}).Finally(func() {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 2"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 2"
		fmt.Println("finally 2")
	})
}
