package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情。
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，我们创建一个新的 promise，将解析的值加上 " Copilot(NewPromise)"
		// In the Then method, we create a new promise, appending " Copilot(NewPromise)" to the resolved value
		return vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve(value.(string) + " Copilot(NewPromise)")
		})
	}, nil).Then(func(value interface{}) interface{} {
		// 在这个 Then 方法中，我们获取前一个 promise 的值，并加上 " !!"
		// In this Then method, we get the value from the previous promise and append " !!"
		return value.(*vl.Promise).GetValue().(string) + " !!"
	}, nil)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(result.GetValue())
}
