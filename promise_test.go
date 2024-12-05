package vowlink

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromise_Then(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return value.(string) + " vowlink", nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink", result.value, "Expected value to be 'Hello, World! vowlink'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Then(nil, func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected reason to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onFulfilled and onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(nil, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Then Chain", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return value.(string) + " vowlink", nil
		}, nil).Then(func(value interface{}) (interface{}, error) {
			return value.(string) + "!", nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink!", result.value, "Expected value to be 'Hello, World! vowlink!'")
	})

	// 当.then中返回的不是promise对象时（包括undefined），p2的状态 一直都是fulfilled，且值为undefined
	t.Run("Then Chain with Rejection", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return value.(string) + " vowlink", nil
		}, nil).Then(func(value interface{}) (interface{}, error) {
			return value.(string) + "!", nil
		}, func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink!", result.value, "Expected value to be 'Hello, World! vowlink!'")
	})

	t.Run("Then return a Promise with resolve", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
				resolve(value.(string)+" vowlink", nil)
			}), nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink", result.value.(*Promise).GetValue(), "Expected value to be 'Hello, World! vowlink'")
	})

	t.Run("Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
				reject(nil, errors.New("Something went wrong"))
			}), nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, Rejected, result.value.(*Promise).state, "Expected value to be Rejected")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).GetReason().Error(), "Expected reason to be 'Something went wrong'")
	})

	t.Run("One Then onRejected after Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value interface{}) (interface{}, error) {
			return NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
				reject(nil, errors.New("Something went wrong"))
			}), nil
		}, nil).Then(nil, func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).reason.Error(), "Expected reason to be 'Something went wrong', Then(nil, func(reason error) error) not work")
	})
}

func TestPromise_Catch(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected value to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Catch(nil)

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.reason.Error(), "Expected value to be 'Something went wrong'")
	})
}

func TestPromise_Finally(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		var finallyCalled bool
		result := p.Finally(func() error {
			finallyCalled = true
			return nil
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
		assert.True(t, finallyCalled, "Expected finally function to be called")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		var finallyCalled bool
		result := p.Finally(func() error {
			finallyCalled = true
			return nil
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Something went wrong", result.reason.Error(), "Expected reason to be 'Something went wrong'")
		assert.True(t, finallyCalled, "Expected finally function to be called")
	})

	t.Run("Nil onFinally", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Finally(nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})
}

func TestMethod_All(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 2 rejected", result.reason.Error(), "Expected reason to be 'Promise 2 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_Any(t *testing.T) {
	t.Run("Any promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, &AggregateError{Errors: []error{errors.New("Promise 1 rejected"), errors.New("Promise 2 rejected"), errors.New("Promise 3 rejected")}}, result.reason, "Expected reason to be an AggregateError")
	})
}

func TestPromise_Race(t *testing.T) {
	t.Run("One promise fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected value to be 'Promise 1 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_AllSettled(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Promise 3", nil)
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{"Promise 1", errors.New("Promise 2 rejected"), "Promise 3"}, result.value, "Expected value to be ['Promise 1', errors.New('Promise 2 rejected'), 'Promise 3']")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []interface{}{errors.New("Promise 1 rejected"), errors.New("Promise 2 rejected"), errors.New("Promise 3 rejected")}, result.value, "Expected value to be [errors.New('Promise 1 rejected'), errors.New('Promise 2 rejected'), errors.New('Promise 3 rejected')]")
	})
}

func TestPromise_MultiCatch(t *testing.T) {
	t.Run("Rejected Multi Catch with New Error", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 3 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong", p.GetReason().Error(), "Expected reason to be 'Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Return Value", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return "Recovered value", nil
		}).Then(func(data interface{}) (interface{}, error) {
			return data, nil
		}, nil)

		assert.Equal(t, "Recovered value", p.GetValue().(string), "Expected value to be 'Recovered value'")
		assert.Nil(t, p.GetReason(), "Expected reason to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Then return error", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return "Recovered value", nil
		}).Then(func(data interface{}) (interface{}, error) {
			return nil, errors.New("Then error: " + data.(string))
		}, nil).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 2 error: Then error: Recovered value", p.GetReason().Error(), "Expected reason to be 'Handled 2 error: Then error: Recovered value'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with New Error and Finally", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		}).Finally(func() error {
			fmt.Println("Finally called")
			return nil
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 3 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong", p.GetReason().Error(), "Expected reason to be 'Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Finally", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (interface{}, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (interface{}, error) {
			return "Recovered value", nil
		}).Finally(func() error {
			fmt.Println("Finally called")
			return nil
		}).Then(func(data interface{}) (interface{}, error) {
			return data, nil
		}, nil)

		assert.Equal(t, "Recovered value", p.GetValue().(string), "Expected value to be 'Recovered value'")
		assert.Nil(t, p.GetReason(), "Expected reason to be nil")
	})
}

func TestPromise_ResolveWithError(t *testing.T) {
	p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		resolve(nil, errors.New("Something went wrong"))
	}).Catch(func(reason error) (interface{}, error) {
		return nil, errors.New("Handled error: " + reason.Error())
	}).Catch(func(reason error) (interface{}, error) {
		return "Recovered value", nil
	}).Finally(func() error {
		fmt.Println("Finally called")
		return nil
	}).Then(func(data interface{}) (interface{}, error) {
		return data, nil
	}, nil)

	assert.Equal(t, "Recovered value", p.GetValue().(string), "Expected value to be 'Recovered value'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
}

func TestPromise_ResolveWithErrorData(t *testing.T) {
	p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		resolve(errors.New("Something went wrong"), nil)
	}).Then(func(data interface{}) (interface{}, error) {
		return data.(error).Error(), nil
	}, func(error) (interface{}, error) {
		return nil, errors.New("Handled error")
	}).Catch(func(reason error) (interface{}, error) {
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil
	})

	assert.Equal(t, "Something went wrong", p.GetValue().(string), "Expected value to be 'Something went wrong'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
}

func TestPromise_RejectWithNil(t *testing.T) {
	p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
		reject("Something went wrong", nil)
	}).Then(func(data interface{}) (interface{}, error) {
		return data, nil
	}, func(error) (interface{}, error) {
		return nil, errors.New("Handled error")
	}).Catch(func(reason error) (interface{}, error) {
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil
	})

	assert.Equal(t, "Something went wrong", p.GetValue().(string), "Expected reason to be 'Something went wrong'")
	assert.Nil(t, p.GetReason(), "Expected value to be nil")
}

func TestPromise_FinallyWithError(t *testing.T) {
	t.Run("Finally with error and resolved", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("Hello, World!", nil)
		}).Finally(func() error {
			return errors.New("Finally error")
		}).Then(func(data interface{}) (interface{}, error) {
			return data.(string) + " vowlink", nil
		}, func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, "Handled error: Finally error", p.GetReason().Error(), "Expected reason to be 'Handled error: Finally error'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")

	})

	t.Run("Finally with error and rejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Finally(func() error {
			return errors.New("Finally error")
		}).Then(func(data interface{}) (interface{}, error) {
			return data.(string) + " vowlink", nil
		}, func(reason error) (interface{}, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, "Handled error: Finally error", p.GetReason().Error(), "Expected reason to be 'Handled error: Finally error'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})
}

func TestNewPromise(t *testing.T) {
	t.Run("nil handler", func(t *testing.T) {
		p := NewPromise(nil)
		assert.Nil(t, p, "Expected nil when handler is nil")
	})

	t.Run("initial state", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			// Empty handler
		})
		assert.Equal(t, Pending, p.state, "Expected initial state to be Pending")
		assert.Nil(t, p.value, "Expected initial value to be nil")
		assert.Nil(t, p.reason, "Expected initial reason to be nil")
	})
}

func TestPromise_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent resolve/reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			go func() {
				resolve("success", nil)
			}()
			go func() {
				reject(nil, errors.New("error"))
			}()
		})

		// Wait for goroutines to complete
		time.Sleep(100 * time.Millisecond)

		// State should be either Fulfilled or Rejected, but not both
		assert.True(t, p.state == Fulfilled || p.state == Rejected,
			"Expected state to be either Fulfilled or Rejected")
		assert.True(t, (p.value == "success" && p.reason == nil) ||
			(p.value == nil && p.reason != nil),
			"Expected either value or reason to be set, not both")
	})
}

func TestPromise_StateImmutability(t *testing.T) {
	t.Run("fulfilled state cannot be changed", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve("first", nil)
			resolve("second", nil)           // should not change state
			reject(nil, errors.New("error")) // should not change state
		})

		assert.Equal(t, Fulfilled, p.state)
		assert.Equal(t, "first", p.value)
		assert.Nil(t, p.reason)
	})

	t.Run("rejected state cannot be changed", func(t *testing.T) {
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			reject(nil, errors.New("first error"))
			reject(nil, errors.New("second error")) // should not change state
			resolve("success", nil)                 // should not change state
		})

		assert.Equal(t, Rejected, p.state)
		assert.Equal(t, "first error", p.reason.Error())
		assert.Nil(t, p.value)
	})
}

func TestPromise_EmptyArray(t *testing.T) {
	t.Run("All with empty array", func(t *testing.T) {
		result := All()
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []interface{}{}, result.value)
	})

	t.Run("Race with empty array", func(t *testing.T) {
		result := Race()
		assert.Equal(t, Fulfilled, result.state)
		assert.Nil(t, result.value)
	})

	t.Run("Any with empty array", func(t *testing.T) {
		result := Any()
		assert.Equal(t, Rejected, result.state)
		assert.IsType(t, &AggregateError{}, result.reason)
	})

	t.Run("AllSettled with empty array", func(t *testing.T) {
		result := AllSettled()
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []interface{}{}, result.value)
	})
}
