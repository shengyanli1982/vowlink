package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {

	// 创建一个新的 Promise
	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 Promise 会立即被拒绝，原因是 "rejected.100"
		// This Promise will be immediately rejected with the reason "rejected.100"
		reject(nil, fmt.Errorf("rejected.100"))

	}).Catch(func(err error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, err

	}).Catch(func(err error) (interface{}, error) {
		// 当上一个 Catch 函数返回错误时，会执行这个 Catch 函数
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 2")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.200")

	}).Catch(func(err error) (interface{}, error) {
		// 当上一个 Catch 函数返回错误时，会执行这个 Catch 函数
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 3")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.300")

	}).Then(func(value interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// 返回一个新的值，将会被下一个 Then 函数接收
		// Return a new value, which will be received by the next Then function
		return fmt.Sprintf("Never be here!! recover value: %v", value), nil

	}, func(err error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 4")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("Should be here.")
	})

	// 打印 Promise 被拒绝的原因
	// Print the reason the Promise was rejected
	fmt.Println("reason: ", result.GetReason())

	// 打印 Promise 的值，但是在这个例子中，Promise 会被拒绝，所以值会是 nil
	// Print the value of the Promise, but in this case, the Promise will be rejected, so the value will be nil
	fmt.Println("value: ", result.GetValue())

}
