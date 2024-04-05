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
		// 这个 Promise 将继续执行，并将错误 "Something went wrong" 作为值传递给下一个 Promise
		// This Promise will continue to execute and pass the error "Something went wrong" as a value to the next Promise
		reject("Something went wrong", nil)

	}).Then(func(data interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// 返回解决的值，将会被下一个 Then 函数接收
		// Return the resolved value, which will be received by the next Then function
		return data, nil

	}, func(error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// 返回一个新的错误 "Handled error"
		// Return a new error "Handled error"
		return nil, errors.New("Handled error")

	}).Catch(func(reason error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 2")

		// 返回一个字符串，表示恢复的值
		// Return a string representing the recovered value
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil

	})

	// 输出 Promise 的拒绝原因
	// Print the rejection reason of the Promise
	fmt.Println("reason: ", result.GetReason())

	// 输出 Promise 的解决值
	// Print the resolution value of the Promise
	fmt.Println("value: ", result.GetValue())

}
