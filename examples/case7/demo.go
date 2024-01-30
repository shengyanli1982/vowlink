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
