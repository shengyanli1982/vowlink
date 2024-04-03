package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情。
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在第一个 Then 方法中，我们将解析的值加上 " vowlink"
		// In the first Then method, we append " vowlink" to the resolved value
		return value.(string) + " vowlink", nil
	}, nil).Then(func(value interface{}) (interface{}, error) {
		// 在第二个 Then 方法中，我们将解析的值加上 " !!"
		// In the second Then method, we append " !!" to the resolved value
		return value.(string) + " !!", nil
	}, nil)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(result.GetValue())
}
