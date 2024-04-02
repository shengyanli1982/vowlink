English | [中文](./README_CN.md)

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
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

**Result**

```bash
$ go run demo.go
hello world Copilot !!
```

It's so easy, right?

Instead of using the same old examples to explain how to use VowLink, let's dive into a real case to demonstrate its usage. By using practical examples, we can better understand how VowLink can help us achieve our goals.

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

**Result**

```bash
$ go run demo.go
Resolve: hello world Copilot !!
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

**Result**

```bash
$ go run demo.go
Resolve: hello world Copilot !!
Rejected: rejected.
```

#### # Case 3

I want to use `finally()` to do something after the promise is resolved or rejected.

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

**Result**

```bash
$ go run demo.go
finally 1
finally 2
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

**Result**

```bash
$ go run demo.go
hello world Copilot(NewPromise) !!
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

**Result**

```bash
$ go run demo.go
>> 0 Promise 1
!! 1 Promise 2 rejected
>> 2 Promise 3
```
