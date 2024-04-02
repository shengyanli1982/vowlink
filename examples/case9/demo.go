package main

import (
	"errors"
	"fmt"

	vl "github.com/shengyanli1982/vowlink"
)

func main() {

	rt := vl.NewPromise(func(resolve func(interface{}), reject func(error)) {

		reject(fmt.Errorf("rejected.100"))

	}).Catch(func(err error) error {

		fmt.Println("> catch 1")
		return err

	}).Catch(func(err error) error {

		fmt.Println("> catch 2")
		return nil

	}).Catch(func(err error) error {

		fmt.Println("> catch 3")
		return errors.New("rejected.200")

	}).Then(func(value interface{}) interface{} {

		fmt.Println("> then 1")

		return "Should be here."

	}, func(err error) error {

		fmt.Println("> catch 4")
		return errors.New("Never be here!!")

	})

	// 这里一定输出 "nil"
	fmt.Println("reason: ", rt.GetReason())

	// 这里一定是 "Should be here."
	fmt.Println("value: ", rt.GetValue())

}
