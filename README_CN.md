[English](./README.md) | 中文

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/vowlink.svg)](https://pkg.go.dev/github.com/shengyanli1982/vowlink)

## 简介

我已经完成了 `VowLink` 的第一个版本，但是几句话语言介绍 `VowLink` 可能有些困难。然而，我相信 `VowLink` 的优点可以通过简单的方式向大家介绍。

在我的开发经验中，我经常遇到代码嵌套过深、逻辑过于复杂的问题。我希望找到一种解决方案，可以避免我陷入代码的中间部分，并简化逻辑。这时我发现了 ES6 中的 Promise 概念，这是一个正确的方向，但并不完美。这激发了我想创造出更好的东西解决问题的动机。

于是，`VowLink` 诞生了。它是一个受 ES6 Promise 启发的 Golang 项目，为函数调用提供了强力的链式调用工具。

## 优势

-   简单易用
-   无第三方依赖
-   支持 `then()`、`catch()`、`finally()`、`all()`、`race()`、`any()`、`allSettled()` 等多种方法

## 安装

```bash
go get github.com/shengyanli1982/vowlink
```

## 快速入门

使用 `VowLink`，您只需几分钟即可开始使用。它非常易于使用，无需进行复杂的设置。

**示例**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink 像一个链条，你可以在链条上添加更多的 then() 来在 promise 解析后做更多的事情。
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error))  {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// 在第一个 Then 方法中，我们将解析的值加上 " vowlink"
		// In the first Then method, we append " vowlink" to the resolved value
		return value.(string) + " vowlink"
	}, nil).Then(func(value interface{}) (interface{}, error) {
		// 在第二个 Then 方法中，我们将解析的值加上 " !!"
		// In the second Then method, we append " !!" to the resolved value
		return value.(string) + " !!"
	}, nil)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(result.GetValue())
}
```

**执行结果**

```bash
$ go run demo.go
hello world vowlink !!
```

这很简单，对吧？

为了说明如何使用 VowLink，我们不再使用相同的旧例子，而是深入一个真实案例来展示它的用法。通过使用实际示例，我们可以更好地理解 VowLink 如何帮助我们实现目标。

### 实例案例

我们的工作中有各种各样的案例，我将展示一些例子。您可以在 `examples` 目录中找到每个案例的代码。例如，案例 1 位于 `examples/case1`。

#### 案例 1

就像在我们的代码中使用 `if` 和 `else` 一样，我们经常希望在条件为真时执行某些操作，在条件为假时执行不同的操作。

**示例**

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

**执行结果**

```bash
$ go run demo.go
Resolve: hello world vowlink !!
Rejected: rejected.
```

#### # 案例 2

我更喜欢在 Golang 中使用类似 JavaScript 风格的代码。我想使用 `then()` 在 promise 解析后执行操作，并使用 `catch()` 处理被拒绝的 promise。

**示例**

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

**执行结果**

```bash
$ go run demo.go
Resolve: hello world vowlink !!
Rejected: rejected.
```

#### # 案例 3

我想使用 `finally()` 在 promise 解析或被拒绝后执行一些操作。

**示例**

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

**执行结果**

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

#### # 案例 4

是的，你可以使用 `then()` 方法返回一个新的 promise。

**示例**

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

**执行结果**

```bash
$ go run demo.go
hello world vowlink(NewPromise) !!
```

#### # 案例 5

如果要在所有 promise 解析后执行某个操作，可以使用 `all()` 方法。

**示例**

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

**执行结果**

```bash
$ go run demo.go
>> 0 Promise 1
>> 1 Promise 2
>> 2 Promise 3
```

#### # 案例 6

我想使用 `race()` 在第一个 promise 解析后执行某个操作。

**示例**

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

**执行结果**

```bash
$ go run demo.go
>> Promise 1
```

#### # 案例 7

我想使用 `any()` 在第一个 promise 解析后执行某个操作。`any()` 类似于 ES6 中的 `race`。然而，`any()` 还会捕获所有 promise 的错误，并在所有 promise 都被拒绝时返回一个 `AggregateError`。

**示例**

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

**执行结果**

```bash
$ go run demo.go
>> Promise 1
!! All promises were rejected: Promise 1 rejected, Promise 2 rejected, Promise 3 rejected
```

#### # 案例 8

要在所有 promise 解析或拒绝后执行操作，可以使用 `allSettled()` 方法。

**示例**

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

**执行结果**

```bash
$ go run demo.go
>> 0 Promise 1
!! 1 Promise 2 rejected
>> 2 Promise 3
```

#### # 案例 9

在创建 Promise 对象后，您可以使用 `reject` 函数触发一个错误。后续的 `catch()` 函数将处理前一个错误并返回一个新的错误，从而创建一个错误调用链。

**示例**

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

**执行结果**

```bash
$ go run demo.go
> catch 1
> catch 2
> catch 3
> catch 4
reason:  Should be here.
value:  <nil>
```

#### # 案例 10

创建一个 Promise 对象后，您可以使用 `reject` 函数触发一个错误。每个后续的 `catch()` 函数将处理前一个错误并返回一个新的错误。如果一个 `catch()` 函数成功从错误中恢复并返回一个正常值（`error` 设置为 `nil`），则所有后续的 `catch()` 函数将不会被执行。相反，`then()` 函数将返回这个值。

**示例**

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

**执行结果**

```bash
$ go run demo.go
> catch 1
> catch 2
> then 1
reason:  <nil>
value:  Should be here. recover value: [error handled]
```

#### # 案例 11

在创建 Promise 对象后，使用 `resolve` 函数来处理正常的响应，但解决的值是一个错误。应该在 Promise 后使用 `then()` 函数来处理错误对象，并将结果存储为值。后续的 `catch()` 函数不会对此错误做出响应。

**示例**

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
```

**执行结果**

```bash
$ go run demo.go
> then 1
reason:  <nil>
value:  Something went wrong
```

#### # 案例 12

在创建 Promise 对象后，你可以使用 `reject` 函数来抛出异常。然而，`reject` 函数并不返回错误，而是返回一个值。只有 `then()` 函数处理 `reject` 传递的值，而后续的 `catch()` 函数会被跳过，不进行任何处理。

**示例**

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
		// 这个 Promise 会立即被拒绝，原因是 "Something went wrong"
		// This Promise will be immediately rejected with the reason "Something went wrong"
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

**执行结果**

```bash
$ go run demo.go
> then 1
reason:  <nil>
value:  Something went wrong
```
