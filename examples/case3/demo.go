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
