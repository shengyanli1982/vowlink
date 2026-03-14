package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

// 定义 main 函数
func main() {

	// 创建一个新的 Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 Promise 将继续执行，并将错误 "Something went wrong" 作为值传递给下一个 Promise
		resolve(errors.New("Something went wrong"), nil)

	}).Then(func(data interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		fmt.Println("> then 1")

		// 返回错误的字符串表示形式
		return data.(error).Error(), nil

	}, func(error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		fmt.Println("> catch 1")

		// 返回一个新的错误 "Handled error"
		return nil, errors.New("Handled error")

	}).Catch(func(reason error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		fmt.Println("> catch 2")

		// 返回一个字符串，表示恢复的值
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil

	})

	// 输出 Promise 的拒绝原因
	fmt.Println("reason: ", result.GetReason())

	// 输出 Promise 的解决值
	fmt.Println("value: ", result.GetValue())

}
