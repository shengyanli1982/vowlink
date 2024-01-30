package vowlink

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromise_Then(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return value.(string) + " Copilot"
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! Copilot", result.value, "Expected value to be 'Hello, World! Copilot'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Something went wrong"))
		})

		result := p.Then(nil, func(reason error) error {
			return errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected reason to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onFulfilled and onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(nil, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Then Chain", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return value.(string) + " Copilot"
		}, nil).Then(func(value interface{}) interface{} {
			return value.(string) + "!"
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! Copilot!", result.value, "Expected value to be 'Hello, World! Copilot!'")
	})

	// 当.then中返回的不是promise对象时（包括undefined），p2的状态 一直都是fulfilled，且值为undefined
	t.Run("Then Chain with Rejection", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return value.(string) + " Copilot"
		}, nil).Then(func(value interface{}) interface{} {
			return value.(string) + "!"
		}, func(reason error) error {
			return errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! Copilot!", result.value, "Expected value to be 'Hello, World! Copilot!'")
	})

	t.Run("Then return a Promise with resolve", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return NewPromise(func(resolve func(interface{}), reject func(error)) {
				resolve(value.(string) + " Copilot")
			})
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! Copilot", result.value.(*Promise).GetValue(), "Expected value to be 'Hello, World! Copilot'")
	})

	t.Run("Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return NewPromise(func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("Something went wrong"))
			})
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, Rejected, result.value.(*Promise).state, "Expected value to be Rejected")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).GetReason().Error(), "Expected reason to be 'Something went wrong'")
	})

	t.Run("One Then onRejected after Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Then(func(value interface{}) interface{} {
			return NewPromise(func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("Something went wrong"))
			})
		}, nil).Then(nil, func(reason error) error {
			return errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).reason.Error(), "Expected reason to be 'Something went wrong', Then(nil, func(reason error) error) not work")
	})
}

func TestPromise_Catch(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Catch(func(reason error) error {
			return errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Something went wrong"))
		})

		result := p.Catch(func(reason error) error {
			return errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected value to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Something went wrong"))
		})

		result := p.Catch(nil)

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.reason.Error(), "Expected value to be 'Something went wrong'")
	})
}

func TestPromise_Finally(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		var finallyCalled bool
		result := p.Finally(func() {
			finallyCalled = true
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
		assert.True(t, finallyCalled, "Expected finally function to be called")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Something went wrong"))
		})

		var finallyCalled bool
		result := p.Finally(func() {
			finallyCalled = true
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Something went wrong", result.reason.Error(), "Expected reason to be 'Something went wrong'")
		assert.True(t, finallyCalled, "Expected finally function to be called")
	})

	t.Run("Nil onFinally", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Hello, World!")
		})

		result := p.Finally(nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})
}

func TestMethod_All(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 2")
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 2 rejected", result.reason.Error(), "Expected reason to be 'Promise 2 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 3 rejected"))
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_Any(t *testing.T) {
	t.Run("Any promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 2")
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 3 rejected"))
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, &AggregateError{Errors: []error{errors.New("Promise 1 rejected"), errors.New("Promise 2 rejected"), errors.New("Promise 3 rejected")}}, result.reason, "Expected reason to be an AggregateError")
	})
}

func TestPromise_Race(t *testing.T) {
	t.Run("One promise fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 2")
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 2")
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected value to be 'Promise 1 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 3 rejected"))
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_AllSettled(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 2")
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 1")
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			resolve("Promise 3")
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", errors.New("Promise 2 rejected"), "Promise 3"}, result.value, "Expected value to be ['Promise 1', errors.New('Promise 2 rejected'), 'Promise 3']")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}), reject func(error)) {
			reject(errors.New("Promise 3 rejected"))
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{errors.New("Promise 1 rejected"), errors.New("Promise 2 rejected"), errors.New("Promise 3 rejected")}, result.value, "Expected value to be [errors.New('Promise 1 rejected'), errors.New('Promise 2 rejected'), errors.New('Promise 3 rejected')]")
	})
}
