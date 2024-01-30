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
