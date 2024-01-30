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
