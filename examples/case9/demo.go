package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {

	// 创建一个新的 Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 Promise 会立即被拒绝，原因是 "rejected.100"
		reject(nil, fmt.Errorf("rejected.100"))

	}).Catch(func(err error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		fmt.Println("> catch 1")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		return nil, err

	}).Catch(func(err error) (interface{}, error) {
		// 当上一个 Catch 函数返回错误时，会执行这个 Catch 函数
		fmt.Println("> catch 2")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		return nil, errors.New("rejected.200")

	}).Catch(func(err error) (interface{}, error) {
		// 当上一个 Catch 函数返回错误时，会执行这个 Catch 函数
		fmt.Println("> catch 3")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		return nil, errors.New("rejected.300")

	}).Then(func(value interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		fmt.Println("> then 1")

		// 返回一个新的值，将会被下一个 Then 函数接收
		return fmt.Sprintf("Never be here!! recover value: %v", value), nil

	}, func(err error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		fmt.Println("> catch 4")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		return nil, errors.New("Should be here.")
	})

	// 打印 Promise 被拒绝的原因
	fmt.Println("reason: ", result.GetReason())

	// 打印 Promise 的值，但是在这个例子中，Promise 会被拒绝，所以值会是 nil
	fmt.Println("value: ", result.GetValue())

}
