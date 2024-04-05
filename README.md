English | [中文](./README_CN.md)

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/vowlink.svg)](https://pkg.go.dev/github.com/shengyanli1982/vowlink)

## Introduction

Introducing `VowLink` in a few words can be challenging. Even though I have finished the first version of `VowLink`, I still struggle to find the perfect way to describe it.

In my development experience, I have often encountered the problem of code being too nested and logic becoming too complex. I wanted to find a solution that would allow me to avoid getting stuck in the middle of the code and simplify the logic. That's when I discovered the promise concept in ES6, which was a step in the right direction but not perfect. This inspired me to create something better.

And thus, `VowLink` was born. It is a Golang project that draws inspiration from ES6 promises and provides a powerful tool for chaining function calls.

## Advantages

-   Simple and easy to use
-   No third-party dependencies
-   Supports various methods such as `then()`, `catch()`, `finally()`, `all()`, `race()`, `any()`, `allSettled()`

## Installation

```bash
go get github.com/shengyanli1982/vowlink
```

## Quick Start

With `VowLink`, you can start using it in just a few minutes. It's incredibly easy to use and requires minimal setup.

**Example**

```go
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
```

**Result**

```bash
$ go run demo.go
hello world vowlink !!
```

It's so easy, right?

Instead of using the same old examples to explain how to use VowLink, let's dive into a real case to demonstrate its usage. By using practical examples, we can better understand how VowLink can help us achieve our goals.

### Core Rules

> [!IMPORTANT]
>
> The `Rules` section is the core of the `VowLink` project. It is important to understand the rules before using `VowLink`.

1. All `Then`, `Catch`, and `Finally` methods of a Promise can return an error. If an error is returned (even if it is the original error), it will be handled by the next `Then` or `Catch` method until the returned error is nil.
2. `vowlink` supports the `resolve` and `reject` methods, which can return data and error, allowing `NewPromise` to make a second choice.
3. The `GetValue` and `GetReason` methods are specific to the Promise and are terminal methods, meaning they do not return a `Promise` object.
4. Although `vowlink` is inspired by `JavaScript Promise`, it is not identical to `JavaScript Promise` as `Golang` has its own differences.

### Study Cases

There are various cases in our work, and I will show you some examples. You can find the code for each case in the `examples` directory. For example, Case 1 is located in `examples/case1`.

#### # Case 1

Just like using `if` and `else` in our code, we often want to perform certain actions if a condition is true, and different actions if the condition is false.

**Example**

```go
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
		// 在 Then 方法中，我们将解析的值加上 " vowlink !!"
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// 如果 promise 被解析，我们将解析的值加上 " vowlink"
		// If the promise is resolved, we append " vowlink" to the resolved value
		return value.(string) + " vowlink", nil
	}, func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// 从 promise 中获取拒绝的原因并打印
	// Get the reason for the rejection from the promise and print it
	fmt.Println("Rejected:", result.GetReason().Error())
}
```

**Result**

```bash
$ go run demo.go
Resolve: hello world vowlink !!
Rejected: rejected.
```

#### # Case 2

I prefer using JavaScript-style code in Golang. I want to use `then()` to perform actions after the promise is resolved and `catch()` to handle rejected promises.

**Example**

```go
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
		// 在 Then 方法中，我们将解析的值加上 " vowlink !!"
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

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
	})

	// 从 promise 中获取拒绝的原因并打印
	// Get the reason for the rejection from the promise and print it
	fmt.Println("Rejected:", result.GetReason().Error())
}
```

**Result**

```bash
$ go run demo.go
Resolve: hello world vowlink !!
Rejected: rejected.
```

#### # Case 3

I want to use `finally()` to do something after the promise is resolved or rejected.

**Example**

```go
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
```

**Result**

```bash
$ go run demo.go
========== finally 1 successfully ==========
finally 1
finally 1 is called. value: hello world vowlink !!, error: <nil>
========== finally 1 error ==========
finally error: 1
finally 1 error, but then is called. value: <nil>, error: Handled error: error in finally 1
========== finally 2 successfully ==========
finally 2
finally 2 is called. value: <nil>, error: rejected.
========== finally 2 error ==========
finally error: 2
finally 2 error, but then is called. value: <nil>, error: Handled error: error in finally 2
```

#### # Case 4

Yes, you can return a new promise using the `then()` method.

**Example**

```go
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
		// 在 Then 方法中，我们创建一个新的 promise，将解析的值加上 " vowlink(NewPromise)"
		// In the Then method, we create a new promise, appending " vowlink(NewPromise)" to the resolved value
		return vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve(value.(string)+" vowlink(NewPromise)", nil)
		}), nil
	}, nil).Then(func(value interface{}) (interface{}, error) {
		// 在这个 Then 方法中，我们获取前一个 promise 的值，并加上 " !!"
		// In this Then method, we get the value from the previous promise and append " !!"
		return value.(*vl.Promise).GetValue().(string) + " !!", nil
	}, nil)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(result.GetValue())
}
```

**Result**

```bash
$ go run demo.go
hello world vowlink(NewPromise) !!
```

#### # Case 5

To perform an action after all promises are resolved, you can use the `all()` method.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// All() 将等待所有的 promise 被解析，并返回一个带有所有值的 promise
	// All() will wait for all promises to be resolved, and return a promise with all the values
	result := vl.All(p1, p2, p3)

	// 从 promise 中获取所有的值并打印
	// Get all the values from the promise and print them
	for i, str := range result.GetValue().([]interface{}) {
		fmt.Println(">>", i, str.(string))
	}
}
```

**Result**

```bash
$ go run demo.go
>> 0 Promise 1
>> 1 Promise 2
>> 2 Promise 3
```

#### # Case 6

I want to use `race()` to perform an action once the first promise is resolved.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// Race() 将等待第一个 promise 被解析，并返回一个带有值的 promise
	// Race() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Race(p1, p2, p3)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(">>", result.GetValue().(string))
}
```

**Result**

```bash
$ go run demo.go
>> Promise 1
```

#### # Case 7

I want to use `any()` to perform an action once the first promise is resolved. `any()` is similar to `race` in ES6. However, `any()` also captures any errors from all promises and returns an `AggregateError` if all promises are rejected.

**Example**

```go
package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// Any() 将等待第一个 promise 被解析，并返回一个带有值的 promise
	// Any() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Any(p1, p2, p3)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(">>", result.GetValue().(string))

	// 创建 3 个 promise
	// Create 3 promises
	p1 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 被拒绝，原因是 "Promise 1 rejected"
		// The first promise is rejected with the reason "Promise 1 rejected"
		reject(nil, errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 被拒绝，原因是 "Promise 3 rejected"
		// The third promise is rejected with the reason "Promise 3 rejected"
		reject(nil, errors.New("Promise 3 rejected"))
	})

	// Any() 将等待所有的 promise 被拒绝，并返回一个带有原因 `AggregateError` 的 promise
	// Any() will wait for all promises to be rejected, and return a promise with the reason `AggregateError`
	result = vl.Any(p1, p2, p3)

	// 从 promise 中获取原因并打印
	// Get the reason from the promise and print it
	fmt.Println("!!", result.GetReason().Error())
}
```

**Result**

```bash
$ go run demo.go
>> Promise 1
!! All promises were rejected: Promise 1 rejected, Promise 2 rejected, Promise 3 rejected
```

#### # Case 8

To perform an action after all promises are resolved or rejected, you can use the `allSettled()` method.

**Example**

```go
package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第一个 promise 直接解析为 "Promise 1"
		// The first promise is directly resolved to "Promise 1"
		resolve("Promise 1", nil)
	})

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// 第三个 promise 直接解析为 "Promise 3"
		// The third promise is directly resolved to "Promise 3"
		resolve("Promise 3", nil)
	})

	// AllSettled() 将等待所有的 promise 被解析或拒绝，并返回一个带有值的 promise
	// AllSettled() will wait for all promises to be resolved or rejected, and return a promise with the value
	result := vl.AllSettled(p1, p2, p3)

	// 从 promise 中获取所有的结果
	// Get all the results from the promise
	for i, r := range result.GetValue().([]interface{}) {
		// 如果结果是一个错误，打印错误信息
		// If the result is an error, print the error message
		if v, ok := r.(error); ok {
			fmt.Println("!!", i, v.Error())
		} else {
			// 否则，打印结果值
			// Otherwise, print the result value
			fmt.Println(">>", i, r.(string))
		}
	}
}
```

**Result**

```bash
$ go run demo.go
>> 0 Promise 1
!! 1 Promise 2 rejected
>> 2 Promise 3
```

#### # Case 9

After creating a Promise object, you can use the `reject` function to trigger an error. Subsequent `catch()` functions will handle the previous error and return a new one, creating a chain of error calls.

```go
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
```

**Result**

```bash
$ go run demo.go
> catch 1
> catch 2
> catch 3
> catch 4
reason:  Should be here.
value:  <nil>
```

#### # Case 10

After creating a Promise object, you can use the `reject` function to trigger an error. Each subsequent `catch()` function will handle the previous error and return a new one. If a `catch()` function successfully recovers from the error and returns a normal value (with `error` set to `nil`), all subsequent `catch()` functions will not be executed. Instead, the `then()` function will return this value.

```go
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

		// 返回一个新的值，将会被下一个 Then 函数接收
		// Return a new value, which will be received by the next Then function
		return "[error handled]", nil

	}).Catch(func(err error) (interface{}, error) {
		// 当上一个 Catch 函数返回错误时，会执行这个 Catch 函数
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 3")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.200")

	}).Then(func(value interface{}) (interface{}, error) {
		// 当 Promise 被解决时，会执行这个 Then 函数
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// 返回一个新的值，将会被下一个 Then 函数接收
		// Return a new value, which will be received by the next Then function
		return fmt.Sprintf("Should be here. recover value: %v", value), nil

	}, func(err error) (interface{}, error) {
		// 当 Promise 被拒绝时，会执行这个 Catch 函数
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 4")

		// 返回一个新的错误，将会被下一个 Catch 函数接收
		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("Never be here!!")

	})

	// 输出 Promise 的拒绝原因，这里一定是 "nil"
	// Print the rejection reason of the Promise, it must be "nil" here
	fmt.Println("reason: ", result.GetReason())

	// 输出 Promise 的解决值，这里一定是 "Should be here."
	// Print the resolution value of the Promise, it must be "Should be here." here
	fmt.Println("value: ", result.GetValue())

}
```

**Result**

```bash
$ go run demo.go
> catch 1
> catch 2
> then 1
reason:  <nil>
value:  Should be here. recover value: [error handled]
```

#### # Case 11

After creating a Promise object, the `resolve` function is used to handle the normal response, but the resolved value is an error. The `then()` function should be used after the Promise to handle the error object and store the result as a value. Subsequent `catch()` functions do not respond to this error.

```go
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
		// 这个 Promise 将继续执行，并将错误 "Something went wrong" 作为值传递给下一个 Promise
		// This Promise will continue to execute and pass the error "Something went wrong" as a value to the next Promise
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
```

**Result**

```bash
$ go run demo.go
> then 1
reason:  <nil>
value:  Something went wrong
```

#### # Case 12

After creating a Promise object, you can use the `reject` function to raise an exception. However, the `reject` function does not return an error, but instead returns a value. Only the `then()` function handles the value passed by `reject`, while subsequent `catch()` functions are skipped without any processing.

```go
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
```

**Result**

```bash
$ go run demo.go
> then 1
reason:  <nil>
value:  Something went wrong
```
