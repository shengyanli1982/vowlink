English | [中文](./README_CN.md)

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/vowlink)](https://goreportcard.com/report/github.com/shengyanli1982/vowlink)
[![Build Status](https://github.com/shengyanli1982/vowlink/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/vowlink/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/vowlink.svg)](https://pkg.go.dev/github.com/shengyanli1982/vowlink)

## Introduction

Ever tried to explain quantum physics to your cat? Well, explaining `VowLink` might be just as challenging! Even after completing its first version, finding the perfect words to describe this elegant solution feels like trying to catch a laser pointer dot - always just out of reach.

As a developer, I've frequently encountered the "callback pyramid of doom" - you know, that moment when your code starts looking like an ASCII art of the Egyptian pyramids. While ES6 Promises were a step in the right direction, I felt there was room for something even better in the Go ecosystem. Thus, `VowLink` was born - a Promise implementation that makes your Go code flow as smoothly as butter on a hot pancake.

## Advantages

-   Simple as pie (and just as delicious to use)
-   Zero dependencies (because sometimes less is more)
-   Full Promise API support including `then()`, `catch()`, `finally()`, `all()`, `race()`, `any()`, and `allSettled()` (it's like having the whole Promise family reunion!)

## Installation

```bash
go get github.com/shengyanli1982/vowlink
```

## Quick Start

Getting started with `VowLink` is easier than making instant noodles. Here's a taste of what it can do:

**Example**

```go
package main

import (
    "fmt"
    vl "github.com/shengyanli1982/vowlink"
)

func main() {
    result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
        resolve("hello world", nil)
    }).Then(func(value interface{}) (interface{}, error) {
        return value.(string) + " vowlink", nil
    }, nil).Then(func(value interface{}) (interface{}, error) {
        return value.(string) + " !!", nil
    }, nil)

    fmt.Println(result.GetValue())
}
```

**Result**

```bash
$ go run demo.go
hello world vowlink !!
```

See? That was smoother than a fresh jar of skippy!

Instead of boring you with theoretical examples that put coffee to shame, let's dive into some real-world scenarios. These practical examples will show you how `VowLink` can make your code dance like nobody's watching.

### Core Rules

> [!IMPORTANT]
>
> Before we jump in, here are the golden rules of `VowLink` - think of them as the "Ten Commandments" of Promise handling in Go.

1. All `Then`, `Catch`, and `Finally` methods can return errors. Like a game of hot potato, these errors will keep bouncing through the chain until someone handles them properly (returns nil).
2. The `resolve` and `reject` methods support both data and error returns, giving `NewPromise` the flexibility of a yoga master.
3. `GetValue` and `GetReason` are terminal methods - they're like the full stop at the end of a sentence. Once called, they don't return a Promise object.
4. While `VowLink` takes inspiration from JavaScript Promises, it's been tailored for Go like a bespoke suit.
5. Don't use goroutines inside `Then()`, `Catch()`, or `Finally()` methods. If you need async operations, wrap the entire Promise in a goroutine - it's like putting the whole party in a separate room.

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
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// If the promise is resolved, we append " vowlink" to the resolved value
		return value.(string) + " vowlink", nil
	}, func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

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
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

	// Get the value from the promise and print it
	fmt.Println("Resolve:", result.GetValue())

	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// If the promise is resolved, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	})

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

	// Print "========== finally 1 successfully ==========" to the console
	fmt.Println("========== finally 1 successfully ==========")

	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 1"
		fmt.Println("finally 1")

		// Return nil indicates that the Finally method was executed successfully
		return nil
	})

	// Use the Printf function to output "finally 1 is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason() as the parameters of the Printf function
	fmt.Printf("finally 1 is called. value: %v, error: %v\n", result.GetValue(), result.GetReason())

	// Print "========== finally 1 error ==========" to the console
	fmt.Println("========== finally 1 error ==========")

	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// The Finally function will be called after the Promise is settled (either resolved or rejected)
		fmt.Println("finally error: 1")

		// Return a new error "error in finally 1"
		return errors.New("error in finally 1")
	}).Then(func(data interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved

		// Return the resolved value, which will be received by the next Then function
		return data.(string) + " vowlink", nil
	}, func(reason error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected

		// Return a new error "Handled error: " plus the error message of the reason
		return nil, errors.New("Handled error: " + reason.Error())
	})

	// Use the Printf function to output "finally 1 error, but then is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason().Error() as the parameters of the Printf function
	fmt.Printf("finally 1 error, but then is called. value: %v, error: %v\n", result.GetValue(), result.GetReason().Error())

	// Print "========== finally 2 successfully ==========" to the console
	fmt.Println("========== finally 2 successfully ==========")

	// This is a rejected promise
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// If the promise is resolved, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// Whether the promise is resolved or rejected, the Finally method will be called and print "finally 2"
		fmt.Println("finally 2")

		// Return nil indicates that the Finally method was executed successfully
		return nil
	})

	// Use the Printf function to output "finally 2 is called. value: %v, error: %v\n" to the console
	// Use result.GetValue() and result.GetReason() as the parameters of the Printf function
	fmt.Printf("finally 2 is called. value: %v, error: %v\n", result.GetValue(), result.GetReason())

	// Print "========== finally 2 error ==========" to the console
	fmt.Println("========== finally 2 error ==========")

	// This is a rejected promise, and an error message is returned in the Finally method
	result = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is rejected with the reason "error"
		reject(nil, fmt.Errorf("error"))
	}).Then(func(value interface{}) (interface{}, error) {
		// If the promise is resolved, we append " vowlink !!" to the resolved value
		return value.(string) + " vowlink !!", nil
	}, nil).Catch(func(err error) (interface{}, error) {
		// If the promise is rejected, we return a new error message "rejected."
		return nil, fmt.Errorf("rejected.")
	}).Finally(func() error {
		// The Finally function will be called after the Promise is settled (either resolved or rejected)

		// Print "finally error: 2" to the console
		fmt.Println("finally error: 2")

		// Return a new error "error in finally 2"
		return errors.New("error in finally 2")

	}).Then(func(data interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved

		// Return the resolved value, which will be received by the next Then function
		return data.(string) + " vowlink", nil
	}, func(reason error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected

		// Return a new error "Handled error: " plus the error message of the reason
		return nil, errors.New("Handled error: " + reason.Error())
	})

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
	// vowlink is like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This promise is directly resolved to "hello world"
		resolve("hello world", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, we create a new promise, appending " vowlink(NewPromise)" to the resolved value
		return vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve(value.(string)+" vowlink(NewPromise)", nil)
		}), nil
	}, nil).Then(func(value interface{}) (interface{}, error) {
		// In this Then method, we get the value from the previous promise and append " !!"
		return value.(*vl.Promise).GetValue().(string) + " !!", nil
	}, nil)

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
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// All() will wait for all promises to be resolved, and return a promise with all the values
	result := vl.All(p1, p2, p3)

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
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// Race() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Race(p1, p2, p3)

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
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The first promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 1" to the resolved value
		return value.(string) + " 1", nil
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The second promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 2" to the resolved value
		return value.(string) + " 2", nil
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The third promise is directly resolved to "Promise"
		resolve("Promise", nil)
	}).Then(func(value interface{}) (interface{}, error) {
		// In the Then method, append " 3" to the resolved value
		return value.(string) + " 3", nil
	}, nil)

	// Any() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Any(p1, p2, p3)

	// Get the value from the promise and print it
	fmt.Println(">>", result.GetValue().(string))

	// Create 3 promises
	p1 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The first promise is rejected with the reason "Promise 1 rejected"
		reject(nil, errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The third promise is rejected with the reason "Promise 3 rejected"
		reject(nil, errors.New("Promise 3 rejected"))
	})

	// Any() will wait for all promises to be rejected, and return a promise with the reason `AggregateError`
	result = vl.Any(p1, p2, p3)

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
	// Create 3 promises
	p1 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The first promise is directly resolved to "Promise 1"
		resolve("Promise 1", nil)
	})

	p2 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The second promise is rejected with the reason "Promise 2 rejected"
		reject(nil, errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// The third promise is directly resolved to "Promise 3"
		resolve("Promise 3", nil)
	})

	// AllSettled() will wait for all promises to be resolved or rejected, and return a promise with the value
	result := vl.AllSettled(p1, p2, p3)

	// Get all the results from the promise
	for i, r := range result.GetValue().([]interface{}) {
		// If the result is an error, print the error message
		if v, ok := r.(error); ok {
			fmt.Println("!!", i, v.Error())
		} else {
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

	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This Promise will be immediately rejected with the reason "rejected.100"
		reject(nil, fmt.Errorf("rejected.100"))

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// Return a new error, which will be received by the next Catch function
		return nil, err

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 2")

		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.200")

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 3")

		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.300")

	}).Then(func(value interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// Return a new value, which will be received by the next Then function
		return fmt.Sprintf("Never be here!! recover value: %v", value), nil

	}, func(err error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 4")

		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("Should be here.")
	})

	// Print the reason the Promise was rejected
	fmt.Println("reason: ", result.GetReason())

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

	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This Promise will be immediately rejected with the reason "rejected.100"
		reject(nil, fmt.Errorf("rejected.100"))

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// Return a new error, which will be received by the next Catch function
		return nil, err

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 2")

		// Return a new value, which will be received by the next Then function
		return "[error handled]", nil

	}).Catch(func(err error) (interface{}, error) {
		// This Catch function will be executed when the previous Catch function returns an error
		fmt.Println("> catch 3")

		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("rejected.200")

	}).Then(func(value interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// Return a new value, which will be received by the next Then function
		return fmt.Sprintf("Should be here. recover value: %v", value), nil

	}, func(err error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 4")

		// Return a new error, which will be received by the next Catch function
		return nil, errors.New("Never be here!!")

	})

	// Print the rejection reason of the Promise, it must be "nil" here
	fmt.Println("reason: ", result.GetReason())

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

// Define the main function
func main() {

	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This Promise will continue to execute and pass the error "Something went wrong" as a value to the next Promise
		resolve(errors.New("Something went wrong"), nil)

	}).Then(func(data interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// Return the string representation of the error
		return data.(error).Error(), nil

	}, func(error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// Return a new error "Handled error"
		return nil, errors.New("Handled error")

	}).Catch(func(reason error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 2")

		// Return a string representing the recovered value
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil

	})

	// Print the rejection reason of the Promise
	fmt.Println("reason: ", result.GetReason())

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

	// Create a new Promise
	result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		// This Promise will continue to execute and pass the error "Something went wrong" as a value to the next Promise
		reject("Something went wrong", nil)

	}).Then(func(data interface{}) (interface{}, error) {
		// This Then function will be executed when the Promise is resolved
		fmt.Println("> then 1")

		// Return the resolved value, which will be received by the next Then function
		return data, nil

	}, func(error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 1")

		// Return a new error "Handled error"
		return nil, errors.New("Handled error")

	}).Catch(func(reason error) (interface{}, error) {
		// This Catch function will be executed when the Promise is rejected
		fmt.Println("> catch 2")

		// Return a string representing the recovered value
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil

	})

	// Print the rejection reason of the Promise
	fmt.Println("reason: ", result.GetReason())

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

#### # Case 13

Here's how to properly handle asynchronous operations with VowLink:

> [!IMPORTANT]
>
> Do not use goroutines (e.g., `go func()`) inside `Then()`, `Catch()`, or `Finally()` methods. If you need asynchronous execution, wrap the entire Promise as a goroutine instead.

```go
package main

import (
	"fmt"
	"time"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// Create a channel to wait for the async operation
	done := make(chan struct{})

	// Launch the async operation
	go func() {
		result := vl.NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			// Simulate some async work
			time.Sleep(1 * time.Second)
			resolve("Async operation completed", nil)
		}).Then(func(value interface{}) (interface{}, error) {
			fmt.Println("> Processing async result")
			return value.(string) + "!", nil
		}, nil)

		fmt.Println("Final result:", result.GetValue())
		close(done)
	}()

	// Wait for the async operation to complete
	<-done
	fmt.Println("Main function completed")
}
```

**Result**

```bash
$ go run demo.go
> Processing async result
Final result: Async operation completed!
Main function completed
```
