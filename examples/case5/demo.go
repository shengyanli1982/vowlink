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
