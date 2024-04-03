package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

// 定义 main 函数
// Define the main function
func main() {

	// 创建一个新的 Promise
	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 Promise 会立即被解决，值是一个错误 "Something went wrong"
		// This Promise will be immediately resolved with a value of an error "Something went wrong"
		resolve(errors.New("Something went wrong"), nil)

	}).Then(func(data interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// 返回错误的字符串表示形式
		// Return the string representation of the error
		return data.(error).Error(), nil

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
