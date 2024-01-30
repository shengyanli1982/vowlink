<div>
	<h1>VowLink</h1>
	<p>VowLink is a Golang project inspired by ES6 promises, offering a powerful tool for your functions chain call.</p>
    <img src="assets/logo.png" alt="logo" width="500px">
</div>

## Introduction

How to introduce VowLink in a few words? Even though now i finish the first version of VowLink, i still can't find a good way to introduce it.

In my development experience, i found that the most common problem is that the code is too nested, and the logic is too complex. I don't want in thinking halt in the middle of the code, and i want to find a way to solve this problem. I found that the promise in ES6 is a good way to solve this problem, but it is not perfect. So i want to make a better promise.

So `VowLink` is born. It is a Golang project inspired by ES6 promises, offering a powerful tool for your functions chain call.

## Advantage

-   Simple and easy to use
-   No third-party dependencies
-   Support `then()`, `catch()`, `finally()`, `all()`, `race()`, `any()`, `allSettled()`

## Installation

```bash
go get github.com/shengyanli1982/vowlink
```

## Quick Start

I want say that `VowLink` is very easy to use, and you can use it in a few minutes.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot"
	}, nil).Then(func(value interface{}) interface{} {
		return value.(string) + " !!"
	}, nil)

	// get the value from the promise
	fmt.Println(result.GetValue())
}
```

**Result**

```bash
$ go run demo.go
hello world Copilot !!
```

It's so easy, right?

In this case i don't want use the same old stuff to explain how to use VowLink, so i will use a real case to show you how to use VowLink. In order to achieve the purpose of this project, i will use same examples to show.

### Study Cases

There are so many cases in our work, and i will show you some of them. Codes are in directory `examples`. eg: case 1 is in `examples/case1`.

#### # Case 1

Like if and else in our code, we want to do something if the condition is true, and do something else if the condition is false.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot !!"
	}, func(error) error {
		return fmt.Errorf("rejected.")
	})

	// get the value from the promise
	fmt.Println("Resolve:", result.GetValue())

	// rejected promise
	result = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot"
	}, func(error) error {
		return fmt.Errorf("rejected.")
	})

	// get the reason from the promise
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

I like JavaScript style code, so i want to use it in Golang. I want to use `then()` to do something after the promise is resolved, and use `catch()` to do something after the promise is rejected.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		return fmt.Errorf("rejected.")
	})

	// get the value from the promise
	fmt.Println("Resolve:", result.GetValue())

	// rejected promise
	result = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		return fmt.Errorf("rejected.")
	})

	// get the reason from the promise
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
	// vowlink like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		return fmt.Errorf("rejected.")
	}).Finally(func() {
		fmt.Println("finally 1")
	})

	// rejected promise
	_ = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(fmt.Errorf("error"))
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " Copilot !!"
	}, nil).Catch(func(error) error {
		return fmt.Errorf("rejected.")
	}).Finally(func() {
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

Can you let a new promise return by `then()`? Yes, you can.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// vowlink like a chain, you can add more then() to the chain to do more things after the promise is resolved.
	result := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("hello world")
	}).Then(func(value interface{}) interface{} {
		return vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve(value.(string) + " Copilot(NewPromise)")
		})
	}, nil).Then(func(value interface{}) interface{} {
		return value.(*vl.Promise).GetValue().(string) + " !!"
	}, nil)

	// get the value from the promise
	fmt.Println(result.GetValue())
}
```

**Result**

```bash
$ go run demo.go
hello world Copilot(NewPromise) !!
```

#### # Case 5

I want to use `all()` to do something after all promises are resolved.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// create 3 promise
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 3"
	}, nil)

	// All() will wait for all promises to be resolved, and return a promise with all the values
	result := vl.All(p1, p2, p3)

	// get the value from the promise
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

I want to use `race()` to do something after the first promise is resolved.

**Example**

```go
package main

import (
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// create 3 promise
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 3"
	}, nil)

	// Race() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Race(p1, p2, p3)

	// get the value from the promise
	fmt.Println(">>", result.GetValue().(string))
}
```

**Result**

```bash
$ go run demo.go
>> Promise 1
```

#### # Case 7

I want to use `any()` to do something after the first promise is resolved. `any()` is the same as `race` in ES6. But `any()` get all promise errors and return AggregateError if all promises are rejected.

**Example**

```go
package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// create 3 promise
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 1"
	}, nil)

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 2"
	}, nil)

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise")
	}).Then(func(value interface{}) interface{} {
		return value.(string) + " 3"
	}, nil)

	// Any() will wait for the first promise to be resolved, and return a promise with the value
	result := vl.Any(p1, p2, p3)

	fmt.Println(">>", result.GetValue().(string))

	// create 3 promise
	p1 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(errors.New("Promise 1 rejected"))
	})

	p2 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(errors.New("Promise 2 rejected"))
	})

	p3 = vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(errors.New("Promise 3 rejected"))
	})

	// Any() will wait for all promises to be rejected, and return a promise with reason `AggregateError`
	result = vl.Any(p1, p2, p3)

	// get the reason from the promise
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

I want to use `allSettled()` to do something after all promises are resolved or rejected.

**Example**

```go
package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {
	// create 3 promise
	p1 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise 1")
	})

	p2 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		reject(errors.New("Promise 2 rejected"))
	})

	p3 := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Promise 3")
	})

	// AllSettled() will wait for all promises to be resolved or rejected, and return a promise with the value
	result := vl.AllSettled(p1, p2, p3)

	// get the all result from the promise
	for i, r := range result.GetValue().([]interface{}) {
		if v, ok := r.(error); ok {
			fmt.Println("!!", i, v.Error())
		} else {
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
