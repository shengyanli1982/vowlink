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
