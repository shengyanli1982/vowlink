package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {

	// 输出 "========== finally 1 successfully ==========" 到控制台
	// Print "========== finally 1 successfully ==========" to the console
	fmt.Println("========== finally 1 successfully ==========")

	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情。
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，我们将解析的值加上 " vowlink !!"
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 1"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 1"
		fmt.Println("finally 1")

		// 返回 nil 表示 Finally 方法执行成功
		// Return nil indicates that the Finally method was executed successfully
		return nil
	})

	// 使用 Printf 函数输出 "finally 1 is called. value: %v, error: %v\n" 到控制台
	// 使用 result.GetValue() 和 result.GetReason() 作为 Printf 函数的参数
	// Use the Printf function to output "finally 1 is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason() as the parameters of the Printf function
	fmt.Printf("finally 1 is called. value: %v, error: %v\n", result.GetValue(), result.GetReason())

	// 输出 "========== finally 1 error ==========" 到控制台
	// Print "========== finally 1 error ==========" to the console
	fmt.Println("========== finally 1 error ==========")

	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，我们将解析的值加上 " vowlink !!"
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// Finally 函数会在 Promise 完成（无论解决还是拒绝）后被调用
		// The Finally function will be called after the Promise is settled (either resolved or rejected)
		fmt.Println("finally error: 1")

		// 返回一个新的错误 "error in finally 1"
		// Return a new error "error in finally 1"
		return errors.New("error in finally 1")
	}).Then(func(data interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved

		// 返回解决的值，将会被下一个 Then 函数接收
		// Return the resolved value, which will be received by the next Then function
		return data.(string) + " vowlink", nil
	}, func(reason error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected

		// 返回一个新的错误 "Handled error: " 加上原因的错误信息
		// Return a new error "Handled error: " plus the error message of the reason
		return nil, errors.New("Handled error: " + reason.Error())
	})

	// 使用 Printf 函数输出 "finally 1 error, but then is called. value: %v, error: %v\n" 到控制台
	// 使用 result.GetValue() 和 result.GetReason().Error() 作为 Printf 函数的参数
	// Use the Printf function to output "finally 1 error, but then is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason().Error() as the parameters of the Printf function
	fmt.Printf("finally 1 error, but then is called. value: %v, error: %v\n", result.GetValue(), result.GetReason().Error())

	// 输出 "========== finally 2 successfully ==========" 到控制台
	// Print "========== finally 2 successfully ==========" to the console
	fmt.Println("========== finally 2 successfully ==========")

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// 如果 promise 被解析，我们将解析的值加上 " vowlink !!"
		// If the promise is resolved, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 2"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 2"
		fmt.Println("finally 2")

		// 返回 nil 表示 Finally 方法执行成功
		// Return nil indicates that the Finally method was executed successfully
		return nil
	})

	// 使用 Printf 函数输出 "finally 2 is called. value: %v, error: %v\n" 到控制台
	// 使用 result.GetValue() 和 result.GetReason() 作为 Printf 函数的参数
	// Use the Printf function to output "finally 2 is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason() as the parameters of the Printf function
	fmt.Printf("finally 2 is called. value: %v, error: %v\n", result.GetValue(), result.GetReason())

	// 输出 "========== finally 2 error ==========" 到控制台
	// Print "========== finally 2 error ==========" to the console
	fmt.Println("========== finally 2 error ==========")

	// 这是一个被拒绝的 promise，Finally 方法中返回了一个错误信息
	// This is a rejected promise, and an error message is returned in the Finally method
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// 如果 promise 被解析，我们将解析的值加上 " vowlink !!"
		// If the promise is resolved, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// Finally 函数会在 Promise 完成（无论解决还是拒绝）后被调用
		// The Finally function will be called after the Promise is settled (either resolved or rejected)

		// 输出 "finally error: 2" 到控制台
		// Print "finally error: 2" to the console
		fmt.Println("finally error: 2")

		// 返回一个新的错误 "error in finally 2"
		// Return a new error "error in finally 2"
		return errors.New("error in finally 2")

	}).Then(func(data interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved

		// 返回解决的值，将会被下一个 Then 函数接收
		// Return the resolved value, which will be received by the next Then function
		return data.(string) + " vowlink", nil
	}, func(reason error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected

		// 返回一个新的错误 "Handled error: " 加上原因的错误信息
		// Return a new error "Handled error: " plus the error message of the reason
		return nil, errors.New("Handled error: " + reason.Error())
	})

	// 使用 Printf 函数输出 "finally 2 error, but then is called. value: %v, error: %v\n" 到控制台
	// 使用 result.GetValue() 和 result.GetReason().Error() 作为 Printf 函数的参数
	// Use the Printf function to output "finally 2 error, but then is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason().Error() as the parameters of the Printf function
	fmt.Printf("finally 2 error, but then is called. value: %v, error: %v\n", result.GetValue(), result.GetReason().Error())
}
