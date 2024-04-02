[English](./README.md) | 中文

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/chebyrash/promise/actions)
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
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在第一个 Then 方法中，我们将解析的值加上 " Copilot"
		// In the first Then method, we append " Copilot" to the resolved value
		return value.(string) + " Copilot"
	}, nil).Then(func(value interface{}) interface{} {
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
hello world Copilot !!
```

这很简单，对吧？

为了说明如何使用 VowLink，我们不再使用相同的旧例子，而是深入一个真实案例来展示它的用法。通过使用实际示例，我们可以更好地理解 VowLink 如何帮助我们实现目标。

### 实例案例

我们的工作中有各种各样的案例，我将展示一些例子。您可以在 `examples` 目录中找到每个案例的代码。例如，案例 1 位于 `examples/case1`。

#### # 案例 1

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
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，我们将解析的值加上 " Copilot !!"
		// In the Then method, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	})

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		// 如果 promise 被解析，我们将解析的值加上 " Copilot"
		// If the promise is resolved, we append " Copilot" to the resolved value
		return value.(string) + " Copilot"
	}, func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	})

	// 从 promise 中获取拒绝的原因并打印
	// Get the reason for the rejection from the promise and print it
	fmt.Println("Rejected:", result.GetReason().Error())
}
```

**执行结果**

```bash
$ go run demo.go
Resolve: hello world Copilot !!
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
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，我们将解析的值加上 " Copilot !!"
		// In the Then method, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	})

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		// 如果 promise 被解析，我们将解析的值加上 " Copilot !!"
		// If the promise is resolved, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	})

	// 从 promise 中获取拒绝的原因并打印
	// Get the reason for the rejection from the promise and print it
	fmt.Println("Rejected:", result.GetReason().Error())
}
```

**执行结果**

```bash
$ go run demo.go
Resolve: hello world Copilot !!
Rejected: rejected.
```

#### # 案例 3

我想使用 `finally()` 在 promise 解析或被拒绝后执行一些操作。

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
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 直接解析为 "hello world"
		// This promise is directly resolved to "hello world"
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，我们将解析的值加上 " Copilot !!"
		// In the Then method, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	}).Finally(func() {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 1"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 1"
		fmt.Println("finally 1")
	})

	// 这是一个被拒绝的 promise
	// This is a rejected promise
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 这个 promise 被拒绝，原因是 "error"
		// This promise is rejected with the reason "error"
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		// 如果 promise 被解析，我们将解析的值加上 " Copilot !!"
		// If the promise is resolved, we append " Copilot !!" to the resolved value
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		// 如果 promise 被拒绝，我们将返回一个新的错误信息 "rejected."
		// If the promise is rejected, we return a new error message "rejected."
		return fmt.Errorf("rejected.")
	}).Finally(func() {
		// 不论 promise 是被解析还是被拒绝，Finally 方法都会被调用，并打印 "finally 2"
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 2"
		fmt.Println("finally 2")
	})
}
```

**执行结果**

```bash
$ go run demo.go
finally 1
finally 2
```

#### # 案例 4

嗯，你可以使用 `then()` 方法返回一个新的 promise。

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
```

**执行结果**

```bash
$ go run demo.go
hello world Copilot(NewPromise) !!
```

#### # 案例 5

要在所有 promise 解析后执行操作，可以使用 `all()` 方法。

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
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3"
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

当第一个 promise 解析后，我想使用 `race()` 执行一个操作。

**Example**

```go
func main() {
	// 创建 3 个 promise
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3"
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

我想使用 `any()` 在第一个 promise 解析后执行一个操作。`any()` 类似于 ES6 中的 `race`。然而，`any()` 还会捕获所有 promise 的错误，并在所有 promise 都被拒绝时返回一个 `AggregateError`。

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
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 直接解析为 "Promise"
		// The first promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 1"
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 直接解析为 "Promise"
		// The second promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 2"
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 直接解析为 "Promise"
		// The third promise is directly resolved to "Promise"
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		// 在 Then 方法中，将解析的值加上 " 3"
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3"
	}, nil)

	// Any() 将等待第一个 promise 被解析，并返回一个带有值的 promise
	// Any() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Any(p1, p2, p3)

	// 从 promise 中获取值并打印
	// Get the value from the promise and print it
	fmt.Println(">>", result.GetValue().(string))

	// 创建 3 个 promise
	// Create 3 promises
	p1 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 被拒绝，原因是 "Promise 1 rejected"
		// The first promise is rejected with the reason "Promise 1 rejected"
		reject(errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 被拒绝，原因是 "Promise 3 rejected"
		// The third promise is rejected with the reason "Promise 3 rejected"
		reject(errors.New("Promise 3 rejected"))
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

要在所有 promise 被解析或拒绝后执行操作，可以使用 `allSettled()` 方法。

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
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第一个 promise 直接解析为 "Promise 1"
		// The first promise is directly resolved to "Promise 1"
		resolve("Promise 1")
	})

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第二个 promise 被拒绝，原因是 "Promise 2 rejected"
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		// 第三个 promise 直接解析为 "Promise 3"
		// The third promise is directly resolved to "Promise 3"
		resolve("Promise 3")
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
