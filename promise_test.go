package vowlink

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
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

	// 当.then中返回的不是promise对象时（包括undefined），p2的状态一直都是fulfilled，且值为undefined
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
		aggErr, ok := result.reason.(*AggregateError)
		assert.True(t, ok, "Expected reason to be an *AggregateError")
		assert.Equal(t, 3, len(aggErr.Errors), "Expected 3 errors in AggregateError")
		assert.Equal(t, "Promise 1 rejected", aggErr.Errors[0].Error())
		assert.Equal(t, "Promise 2 rejected", aggErr.Errors[1].Error())
		assert.Equal(t, "Promise 3 rejected", aggErr.Errors[2].Error())
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

		assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
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

		assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
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

	assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
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

	assert.Equal(t, "Recovered value: Handled error", p.GetValue().(string), "Expected value to be 'Recovered value: Handled error'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
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

		assert.Equal(t, "Handled error: Something went wrong\nFinally error", p.GetReason().Error(), "Expected reason to contain both original and cleanup errors")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})
}

func TestNewPromise(t *testing.T) {
	t.Run("nil handler", func(t *testing.T) {
		p := NewPromise(nil)
		assert.NotNil(t, p, "Expected rejected Promise when handler is nil")
		assert.Equal(t, Rejected, p.state, "Expected state to be Rejected")
		assert.NotNil(t, p.GetReason(), "Expected reason to be non-nil")
		assert.Equal(t, "promise handler cannot be nil", p.GetReason().Error())
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
		var wg sync.WaitGroup
		wg.Add(2)

		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			go func() {
				defer wg.Done()
				resolve("success", nil)
			}()
			go func() {
				defer wg.Done()
				reject(nil, errors.New("error"))
			}()
		})

		wg.Wait()

		// State should be either Fulfilled or Rejected, but not both
		state := p.getState()
		value := p.GetValue()
		reason := p.GetReason()

		assert.True(t, state == Fulfilled || state == Rejected,
			"Expected state to be either Fulfilled or Rejected")
		assert.True(t, (value == "success" && reason == nil) ||
			(value == nil && reason != nil),
			"Expected either value or reason to be set, not both")
	})
}

func TestPromise_NoGoroutineLeak(t *testing.T) {
	t.Run("concurrent settle should not leak goroutines", func(t *testing.T) {
		const iterations = 2000

		before := runtime.NumGoroutine()
		for i := 0; i < iterations; i++ {
			var wg sync.WaitGroup
			wg.Add(2)

			_ = NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
				go func() {
					defer wg.Done()
					resolve("ok", nil)
				}()
				go func() {
					defer wg.Done()
					reject(nil, errors.New("error"))
				}()
			})

			wg.Wait()
		}

		// 在限定时间内轮询，降低慢机上的抖动误报。
		deadline := time.Now().Add(500 * time.Millisecond)
		for {
			runtime.GC()
			after := runtime.NumGoroutine()
			if after <= before+4 {
				return
			}
			if time.Now().After(deadline) {
				assert.LessOrEqual(t, after, before+4, "possible goroutine leak detected: before=%d after=%d", before, after)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
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

func TestPromise_AsyncThen(t *testing.T) {
	t.Run("async resolve with Then chain", func(t *testing.T) {
		done := make(chan struct{})
		var result string
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				resolve("async", nil)
			}()
		})
		p.Then(func(value interface{}) (interface{}, error) {
			result = value.(string) + " done"
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for async Then")
		}
		assert.Equal(t, "async done", result)
	})

	t.Run("async reject with Catch chain", func(t *testing.T) {
		done := make(chan struct{})
		var errMsg string
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				reject(nil, errors.New("async error"))
			}()
		})
		p.Catch(func(reason error) (interface{}, error) {
			errMsg = reason.Error()
			close(done)
			return nil, nil
		})
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for async Catch")
		}
		assert.Equal(t, "async error", errMsg)
	})

	t.Run("multiple Then on pending promise", func(t *testing.T) {
		done1 := make(chan struct{})
		done2 := make(chan struct{})
		var r1, r2 int
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve(42, nil)
			}()
		})
		p.Then(func(value interface{}) (interface{}, error) {
			r1 = value.(int) * 2
			close(done1)
			return nil, nil
		}, nil)
		p.Then(func(value interface{}) (interface{}, error) {
			r2 = value.(int) * 3
			close(done2)
			return nil, nil
		}, nil)
		select {
		case <-done1:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout on done1")
		}
		select {
		case <-done2:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout on done2")
		}
		assert.Equal(t, 84, r1)
		assert.Equal(t, 126, r2)
	})
}

func TestPromise_AsyncAll(t *testing.T) {
	t.Run("async resolve collection with correct order", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("first", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve("second", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("third", nil)
			}()
		})
		result := All(p1, p2, p3)
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for All async")
		}
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]interface{})
		assert.Equal(t, 3, len(values))
		assert.Equal(t, "first", values[0])
		assert.Equal(t, "second", values[1])
		assert.Equal(t, "third", values[2])
	})

	t.Run("one async reject fails All", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("ok", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("fast fail"))
			}()
		})
		result := All(p1, p2)
		result.Then(nil, func(reason error) (interface{}, error) {
			close(done)
			return nil, nil
		})
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for All reject")
		}
		assert.Equal(t, Rejected, result.state)
		assert.Equal(t, "fast fail", result.reason.Error())
	})
}

func TestPromise_AsyncRace(t *testing.T) {
	t.Run("fastest resolve wins", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				resolve("slow", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				resolve("fast", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(300 * time.Millisecond)
				resolve("slowest", nil)
			}()
		})
		result := Race(p1, p2, p3)
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for Race")
		}
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, "fast", result.value)
	})

	t.Run("fastest reject wins", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				resolve("slow", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				reject(nil, errors.New("fast error"))
			}()
		})
		result := Race(p1, p2)
		result.Then(nil, func(reason error) (interface{}, error) {
			close(done)
			return nil, nil
		})
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for Race reject")
		}
		assert.Equal(t, Rejected, result.state)
		assert.Equal(t, "fast error", result.reason.Error())
	})
}

func TestPromise_AsyncAny(t *testing.T) {
	t.Run("first resolve wins among mixed", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				reject(nil, errors.New("err1"))
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve("winner", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("loser", nil)
			}()
		})
		result := Any(p1, p2, p3)
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for Any")
		}
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, "winner", result.value)
	})

	t.Run("all async reject with AggregateError", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("async err1"))
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				reject(nil, errors.New("async err2"))
			}()
		})
		result := Any(p1, p2)
		result.Then(nil, func(reason error) (interface{}, error) {
			close(done)
			return nil, nil
		})
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for Any all reject")
		}
		assert.Equal(t, Rejected, result.state)
		aggErr, ok := result.reason.(*AggregateError)
		assert.True(t, ok)
		assert.Equal(t, 2, len(aggErr.Errors))
	})
}

func TestPromise_AsyncAllSettled(t *testing.T) {
	t.Run("mixed async resolve and reject", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("a", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("err b"))
			}()
		})
		p3 := NewPromise(func(resolve, reject func(interface{}, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("c", nil)
			}()
		})
		result := AllSettled(p1, p2, p3)
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for AllSettled")
		}
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]interface{})
		assert.Equal(t, 3, len(values))
		assert.Equal(t, "a", values[0])
		assert.EqualError(t, errors.New("err b"), values[1].(error).Error())
		assert.Equal(t, "c", values[2])
	})
}

func TestPromise_NilPromise(t *testing.T) {
	t.Run("All with nil in middle", func(t *testing.T) {
		p1 := NewPromise(func(resolve, reject func(interface{}, error)) {
			resolve("a", nil)
		})
		p2 := NewPromise(func(resolve, reject func(interface{}, error)) {
			resolve("b", nil)
		})
		result := All(p1, nil, p2)
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]interface{})
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "a", values[0])
		assert.Equal(t, "b", values[1])
	})

	t.Run("All with only nil", func(t *testing.T) {
		result := All(nil)
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []interface{}{}, result.value)
	})

	t.Run("AllSettled with nil nil", func(t *testing.T) {
		result := AllSettled(nil, nil)
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []interface{}{}, result.value)
	})

	t.Run("Any with nil nil", func(t *testing.T) {
		result := Any(nil, nil)
		assert.Equal(t, Rejected, result.state)
		assert.IsType(t, &AggregateError{}, result.reason)
	})

	t.Run("Race with nil", func(t *testing.T) {
		result := Race(nil)
		assert.Equal(t, Fulfilled, result.state)
		assert.Nil(t, result.value)
	})
}

func TestPromise_ConcurrentAll(t *testing.T) {
	t.Run("large number of concurrent promises", func(t *testing.T) {
		const n = 1000
		promises := make([]*Promise, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			idx := i
			promises[i] = NewPromise(func(resolve, reject func(interface{}, error)) {
				go func() {
					defer wg.Done()
					time.Sleep(time.Duration(idx%10) * time.Millisecond)
					resolve(idx, nil)
				}()
			})
		}
		result := All(promises...)
		done := make(chan struct{})
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, func(reason error) (interface{}, error) {
			close(done)
			return nil, nil
		})
		wg.Wait()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for large All")
		}
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]interface{})
		assert.Equal(t, n, len(values))
		for i := 0; i < n; i++ {
			assert.Equal(t, i, values[i].(int))
		}
	})

	t.Run("concurrent All with mixed timing", func(t *testing.T) {
		const n = 200
		promises := make([]*Promise, n)
		for i := 0; i < n; i++ {
			idx := i
			promises[i] = NewPromise(func(resolve, reject func(interface{}, error)) {
				go func() {
					time.Sleep(time.Duration(idx%5) * time.Millisecond)
					if idx%7 == 0 {
						reject(nil, errors.New("err"))
					} else {
						resolve(idx, nil)
					}
				}()
			})
		}
		result := All(promises...)
		done := make(chan struct{})
		result.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return nil, nil
		}, func(reason error) (interface{}, error) {
			close(done)
			return nil, nil
		})
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent mixed All")
		}
		assert.True(t, result.state == Fulfilled || result.state == Rejected)
	})
}

func TestPromise_ExecutorPanic(t *testing.T) {
	t.Run("executor panic recovers to rejected", func(t *testing.T) {
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			panic("executor boom")
		})
		assert.Equal(t, Rejected, p.state)
		assert.NotNil(t, p.GetReason())
		assert.Contains(t, p.GetReason().Error(), "promise executor panic")
		assert.Contains(t, p.GetReason().Error(), "executor boom")
	})

	t.Run("executor panic does not block Then chain", func(t *testing.T) {
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			panic("executor boom")
		}).Catch(func(reason error) (interface{}, error) {
			return "recovered from: " + reason.Error(), nil
		})
		assert.Equal(t, Fulfilled, p.state)
		assert.Contains(t, p.GetValue().(string), "recovered from")
	})
}

func TestPromise_SubscriberPanicProtection(t *testing.T) {
	t.Run("panicking subscriber does not block others", func(t *testing.T) {
		var asyncResolve func(interface{}, error)
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			asyncResolve = resolve
		})

		// First subscriber will panic
		result1 := p.Then(func(value interface{}) (interface{}, error) {
			panic("subscriber 1 panic")
		}, nil)

		// Second subscriber should still be called
		done := make(chan struct{})
		var result2Value interface{}
		_ = p.Then(func(value interface{}) (interface{}, error) {
			result2Value = value
			close(done)
			return value, nil
		}, nil)

		// Resolve the promise
		asyncResolve("hello", nil)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout: second subscriber was never called")
		}

		// result1's downstream should be rejected due to panic
		assert.Equal(t, Rejected, result1.state)
		assert.NotNil(t, result1.GetReason())
		assert.Contains(t, result1.GetReason().Error(), "subscriber callback panic")

		// result2's value should be set normally
		assert.Equal(t, "hello", result2Value)
	})
}

func TestPromise_SettlePanicProtection(t *testing.T) {
	t.Run("multiple subscribers with mixed panic - first two panic, third succeeds (fulfilled)", func(t *testing.T) {
		var asyncResolve func(interface{}, error)
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			asyncResolve = resolve
		})

		// First subscriber panics
		result1 := p.Then(func(value interface{}) (interface{}, error) {
			panic("subscriber 1 boom")
		}, nil)

		// Second subscriber panics
		result2 := p.Then(func(value interface{}) (interface{}, error) {
			panic("subscriber 2 boom")
		}, nil)

		// Third subscriber succeeds
		done := make(chan struct{})
		result3 := p.Then(func(value interface{}) (interface{}, error) {
			close(done)
			return value, nil
		}, nil)

		asyncResolve("success", nil)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout: third subscriber was never called")
		}

		// result1 downstream: Rejected (panic recovered)
		assert.Equal(t, Rejected, result1.state)
		assert.NotNil(t, result1.GetReason())
		assert.Contains(t, result1.GetReason().Error(), "subscriber callback panic")
		assert.Contains(t, result1.GetReason().Error(), "subscriber 1 boom")

		// result2 downstream: Rejected (panic recovered)
		assert.Equal(t, Rejected, result2.state)
		assert.NotNil(t, result2.GetReason())
		assert.Contains(t, result2.GetReason().Error(), "subscriber callback panic")
		assert.Contains(t, result2.GetReason().Error(), "subscriber 2 boom")

		// result3 downstream: Fulfilled with correct value
		assert.Equal(t, Fulfilled, result3.state)
		assert.Nil(t, result3.GetReason())
		assert.Equal(t, "success", result3.GetValue())
	})

	t.Run("single panicking subscriber in rejected path does not block others", func(t *testing.T) {
		var asyncReject func(interface{}, error)
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			asyncReject = reject
		})

		// First subscriber panics
		result1 := p.Then(nil, func(reason error) (interface{}, error) {
			panic("reject handler panic")
		})

		// Second subscriber handles rejection normally
		done := make(chan struct{})
		result2 := p.Then(nil, func(reason error) (interface{}, error) {
			close(done)
			return nil, errors.New("handled: " + reason.Error())
		})

		asyncReject(nil, errors.New("original rejection"))

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout: second subscriber was never called")
		}

		// result1 downstream: Rejected (panic recovered)
		assert.Equal(t, Rejected, result1.state)
		assert.NotNil(t, result1.GetReason())
		assert.Contains(t, result1.GetReason().Error(), "subscriber callback panic")
		assert.Contains(t, result1.GetReason().Error(), "reject handler panic")

		// result2 downstream: Rejected with handled error
		assert.Equal(t, Rejected, result2.state)
		assert.NotNil(t, result2.GetReason())
		assert.Equal(t, "handled: original rejection", result2.GetReason().Error())
	})
}

func TestAggregateError_NilError(t *testing.T) {
	t.Run("nil errors rendered as <nil>", func(t *testing.T) {
		ae := &AggregateError{
			Errors: []error{errors.New("err1"), nil, errors.New("err2")},
		}
		assert.Equal(t, "All promises were rejected: err1, <nil>, err2", ae.Error())
	})
}
// ---------- Round 3: 新增边界场景测试 ----------

// TestPromise_MultiSubscriberPanicLifecycle
// 场景: 一个 Promise 有 5 个 Then subscriber，其中第 1、3、5 个 subscriber 回调 panic，第 2、4 个正常。
// 验证:
//   - 第 1、3、5 个下游 Promise 状态为 Rejected，reason 包含 "subscriber callback panic"
//   - 第 2、4 个下游 Promise 状态为 Fulfilled，值正确传递
//   - 验证所有 subscriber 都被调用了（不会因为前面的 panic 而遗漏）
func TestPromise_MultiSubscriberPanicLifecycle(t *testing.T) {
	t.Run("mixed panic and normal subscribers", func(t *testing.T) {
		// Arrange: 创建 Pending 状态的 Promise，保留 resolve 引用用于手动触发
		var asyncResolve func(interface{}, error)
		source := NewPromise(func(resolve, reject func(interface{}, error)) {
			asyncResolve = resolve
		})

		// 用 channel 追踪每个 subscriber 是否被调用
		const n = 5
		called := make([]chan struct{}, n)
		for i := 0; i < n; i++ {
			called[i] = make(chan struct{})
		}

		// 5 个下游 Promise — 注意必须用局部变量捕获索引，避免闭包共享同一变量
		downstream := make([]*Promise, n)

		// subscriber #1: panic
		func() {
			idx := 0
			downstream[idx] = source.Then(func(value interface{}) (interface{}, error) {
				close(called[idx])
				panic("boom from subscriber 1")
			}, nil)
		}()

		// subscriber #2: 正常
		func() {
			idx := 1
			downstream[idx] = source.Then(func(value interface{}) (interface{}, error) {
				close(called[idx])
				return fmt.Sprintf("sub2 got %v", value), nil
			}, nil)
		}()

		// subscriber #3: panic
		func() {
			idx := 2
			downstream[idx] = source.Then(func(value interface{}) (interface{}, error) {
				close(called[idx])
				panic("boom from subscriber 3")
			}, nil)
		}()

		// subscriber #4: 正常
		func() {
			idx := 3
			downstream[idx] = source.Then(func(value interface{}) (interface{}, error) {
				close(called[idx])
				return fmt.Sprintf("sub4 got %v", value), nil
			}, nil)
		}()

		// subscriber #5: panic
		func() {
			idx := 4
			downstream[idx] = source.Then(func(value interface{}) (interface{}, error) {
				close(called[idx])
				panic("boom from subscriber 5")
			}, nil)
		}()

		// Act: resolve 源 Promise，触发所有 subscriber
		asyncResolve("hello", nil)

		// Assert: 验证所有 5 个 subscriber 都被调用
		for idx := 0; idx < n; idx++ {
			select {
			case <-called[idx]:
				// subscriber idx 被调用了
			case <-time.After(2 * time.Second):
				t.Fatalf("timeout: subscriber #%d was never called", idx+1)
			}
		}

		// Assert: 第 1、3、5 个下游应为 Rejected，reason 包含 "subscriber callback panic"
		panicIndices := []int{0, 2, 4}
		panicMessages := []string{"boom from subscriber 1", "boom from subscriber 3", "boom from subscriber 5"}
		for j, idx := range panicIndices {
			assert.Equal(t, Rejected, downstream[idx].getState(),
				"subscriber #%d downstream should be Rejected", idx+1)
			assert.NotNil(t, downstream[idx].GetReason(),
				"subscriber #%d downstream reason should not be nil", idx+1)
			assert.Contains(t, downstream[idx].GetReason().Error(), "subscriber callback panic",
				"subscriber #%d reason should contain 'subscriber callback panic'", idx+1)
			assert.Contains(t, downstream[idx].GetReason().Error(), panicMessages[j],
				"subscriber #%d reason should contain original panic message", idx+1)
		}

		// Assert: 第 2、4 个下游应为 Fulfilled，值正确传递
		normalIndices := []int{1, 3}
		normalValues := []string{"sub2 got hello", "sub4 got hello"}
		for j, idx := range normalIndices {
			assert.Equal(t, Fulfilled, downstream[idx].getState(),
				"subscriber #%d downstream should be Fulfilled", idx+1)
			assert.Nil(t, downstream[idx].GetReason(),
				"subscriber #%d downstream reason should be nil", idx+1)
			assert.Equal(t, normalValues[j], downstream[idx].GetValue(),
				"subscriber #%d downstream value mismatch", idx+1)
		}
	})
}

// TestPromise_DeepChainStackSafety
// 场景: 创建 5000 层同步 Then 链（每个 Promise 已在 executor 中 resolve）。
// 验证:
//   - 不会 panic / 栈溢出
//   - 最终值正确传递
func TestPromise_DeepChainStackSafety(t *testing.T) {
	t.Run("5000 layer synchronous Then chain", func(t *testing.T) {
		// Arrange: 起始值 = 1
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve(1, nil)
		})

		// Act: 5000 层同步 Then 链，每层将值 +1
		const layers = 5000
		assert.NotPanics(t, func() {
			for i := 0; i < layers; i++ {
				p = p.Then(func(value interface{}) (interface{}, error) {
					return value.(int) + 1, nil
				}, nil)
			}
		}, "building 5000-layer Then chain should not panic")

		// Assert: 最终值 = 1 + 5000 = 5001
		assert.Equal(t, Fulfilled, p.getState(),
			"final Promise should be Fulfilled after 5000 layers")
		assert.Equal(t, 5001, p.GetValue(),
			"final value should be 5001 (1 initial + 5000 increments)")
		assert.Nil(t, p.GetReason(),
			"final Promise should have no rejection reason")
	})

	t.Run("5000 layer synchronous Then chain with rejection propagation", func(t *testing.T) {
		// 额外验证: 深层链中间某层 reject 后，后续层不再执行 onSuccess
		p := NewPromise(func(resolve func(interface{}, error), reject func(interface{}, error)) {
			resolve(1, nil)
		})

		const rejectAt = 2500
		successCallCount := 0

		assert.NotPanics(t, func() {
			for i := 0; i < 5000; i++ {
				layer := i
				p = p.Then(func(value interface{}) (interface{}, error) {
					successCallCount++
					if layer == rejectAt {
						return nil, errors.New("mid-chain rejection")
					}
					return value.(int) + 1, nil
				}, nil)
			}
		}, "building deep chain with mid-rejection should not panic")

		// Assert: onSuccess 只在层 0~rejectAt 被调用（共 rejectAt+1 次）
		// rejectAt 层返回 (nil, error) → dispatchCallback 调 reject → 下游 Rejected
		// rejectAt+1 层开始 state=Rejected，走 defaultErrorHandler(nil, reason) → 继续 reject 传播
		// 所以 successCallCount 应恰好等于 rejectAt+1
		assert.Equal(t, rejectAt+1, successCallCount,
			"success handler should be called exactly rejectAt+1 times (layers 0..%d)", rejectAt)

		// 最终 Promise 应为 Rejected，携带 mid-chain rejection 错误
		state, _, reason := p.snapshot()
		assert.Equal(t, Rejected, state,
			"final Promise should be Rejected after mid-chain error")
		assert.NotNil(t, reason,
			"final Promise should carry the rejection reason")
		assert.Equal(t, "mid-chain rejection", reason.Error())
	})
}

// TestPromise_NilPromiseChainCalls
// 场景: NewPromise(nil) 返回一个 Rejected Promise (Round 1 已修复)，
//       验证其 Then/Catch/Finally 链式调用正常工作。
// 验证:
//   - NewPromise(nil).Then(...) 不 panic
//   - NewPromise(nil).Catch(...) 可以捕获 "promise handler cannot be nil" 错误
//   - NewPromise(nil).Finally(...) 正常执行
//   - 链式调用 NewPromise(nil).Catch(recover).Then(use) 正常工作
func TestPromise_NilPromiseChainCalls(t *testing.T) {
	t.Run("Then on nil handler Promise", func(t *testing.T) {
		// NewPromise(nil) 返回 {state: Rejected, reason: "promise handler cannot be nil"}
		// .Then(success, error) 在 Rejected 状态下会调用 errorHandler
		// 因为 errorHandler == nil → 使用 defaultErrorHandler → 直接传递 error → reject
		nilPromise := NewPromise(nil)
		assert.Equal(t, Rejected, nilPromise.getState(), "precondition: nil handler Promise should be Rejected")

		// Then with explicit error handler: 捕获 error 并返回恢复值
		assert.NotPanics(t, func() {
			result := nilPromise.Then(
				func(value interface{}) (interface{}, error) {
					t.Fatal("success handler should not be called on Rejected Promise")
					return nil, nil
				},
				func(reason error) (interface{}, error) {
					return "recovered via Then: " + reason.Error(), nil
				},
			)
			assert.Equal(t, Fulfilled, result.getState(), "recovered Then result should be Fulfilled")
			assert.Equal(t, "recovered via Then: promise handler cannot be nil", result.GetValue())
		}, "Then on nil handler Promise should not panic")

		// Then with nil handlers: defaultErrorHandler 传播 rejection
		assert.NotPanics(t, func() {
			result := nilPromise.Then(nil, nil)
			assert.Equal(t, Rejected, result.getState(),
				"Then(nil,nil) on Rejected should propagate Rejected")
			assert.Equal(t, "promise handler cannot be nil", result.GetReason().Error())
		}, "Then(nil,nil) on nil handler Promise should not panic")
	})

	t.Run("Catch on nil handler Promise", func(t *testing.T) {
		nilPromise := NewPromise(nil)

		// Catch 可以捕获 "promise handler cannot be nil"
		assert.NotPanics(t, func() {
			result := nilPromise.Catch(func(reason error) (interface{}, error) {
				return "caught: " + reason.Error(), nil
			})
			assert.Equal(t, Fulfilled, result.getState(),
				"Catch on nil handler Promise should recover to Fulfilled")
			assert.Equal(t, "caught: promise handler cannot be nil", result.GetValue(),
				"Catch handler should receive the nil handler error message")
		}, "Catch on nil handler Promise should not panic")

		// Catch with nil handler: defaultErrorHandler 传播
		assert.NotPanics(t, func() {
			result := nilPromise.Catch(nil)
			assert.Equal(t, Rejected, result.getState(),
				"Catch(nil) should propagate Rejected state")
			assert.Equal(t, "promise handler cannot be nil", result.GetReason().Error())
		}, "Catch(nil) on nil handler Promise should not panic")
	})

	t.Run("Finally on nil handler Promise", func(t *testing.T) {
		nilPromise := NewPromise(nil)

		// Finally 应正常执行 cleanup，且不 panic
		assert.NotPanics(t, func() {
			finallyCalled := false
			result := nilPromise.Finally(func() error {
				finallyCalled = true
				return nil
			})
			assert.Equal(t, Rejected, result.getState(),
				"Finally on Rejected should remain Rejected (cleanup returns nil)")
			assert.True(t, finallyCalled, "Finally cleanup should be called on nil handler Promise")
			assert.Equal(t, "promise handler cannot be nil", result.GetReason().Error(),
				"Finally should preserve the original rejection reason")
		}, "Finally on nil handler Promise should not panic")

		// Finally with nil handler
		assert.NotPanics(t, func() {
			result := nilPromise.Finally(nil)
			assert.Equal(t, Rejected, result.getState(),
				"Finally(nil) on nil handler Promise should remain Rejected")
		}, "Finally(nil) on nil handler Promise should not panic")

		// Finally with cleanup error: errors.Join(reason, cleanupErr)
		assert.NotPanics(t, func() {
			result := nilPromise.Finally(func() error {
				return errors.New("cleanup failed")
			})
			assert.Equal(t, Rejected, result.getState(),
				"Finally with cleanup error should remain Rejected")
			assert.NotNil(t, result.GetReason(),
				"should carry joined error")
			assert.Contains(t, result.GetReason().Error(), "promise handler cannot be nil",
				"joined error should contain original reason")
			assert.Contains(t, result.GetReason().Error(), "cleanup failed",
				"joined error should contain cleanup error")
		}, "Finally with cleanup error should not panic")
	})

	t.Run("full chain on nil handler Promise", func(t *testing.T) {
		// 完整链: NewPromise(nil).Catch(recover).Finally(log).Then(use)
		nilPromise := NewPromise(nil)

		assert.NotPanics(t, func() {
			finallyCalled := false

			result := nilPromise.
				Catch(func(reason error) (interface{}, error) {
					// 从 nil handler 错误中恢复
					return "recovered from nil handler", nil
				}).
				Finally(func() error {
					finallyCalled = true
					return nil
				}).
				Then(func(value interface{}) (interface{}, error) {
					return value.(string) + " and processed", nil
				}, nil)

			assert.True(t, finallyCalled, "Finally should be called in the chain")
			assert.Equal(t, Fulfilled, result.getState(),
				"full chain result should be Fulfilled")
			assert.Equal(t, "recovered from nil handler and processed", result.GetValue(),
				"full chain should pass recovered value through Finally and Then")
		}, "full Catch→Finally→Then chain on nil handler Promise should not panic")

		// 链式调用不恢复: Catch 再抛出 error
		assert.NotPanics(t, func() {
			result := nilPromise.
				Catch(func(reason error) (interface{}, error) {
					return nil, errors.New("new error from Catch")
				}).
				Finally(func() error {
					return nil
				}).
				Then(func(value interface{}) (interface{}, error) {
					t.Fatal("success handler should not be called on error propagation path")
					return nil, nil
				}, func(reason error) (interface{}, error) {
					return "handled: " + reason.Error(), nil
				})

			assert.Equal(t, Fulfilled, result.getState(),
				"final handler should recover the chain to Fulfilled")
			assert.Equal(t, "handled: new error from Catch", result.GetValue())
		}, "Catch→Finally→Then chain with error re-throw should not panic")
	})
}

func TestPromise_ThenPendingToSettledRace(t *testing.T) {
	t.Run("concurrent Then registration and settle", func(t *testing.T) {
		// Exercise the double-check path in Then() where the promise
		// transitions from Pending to settled between snapshot() and Lock().
		const iterations = 200
		for i := 0; i < iterations; i++ {
			var asyncResolve func(interface{}, error)
			p := NewPromise(func(resolve, reject func(interface{}, error)) {
				asyncResolve = resolve
			})

			// Use channel to signal completion from Then callback
			done := make(chan string, 1)
			p.Then(func(v interface{}) (interface{}, error) {
				done <- v.(string) + "!"
				return v, nil
			}, nil)

			// Resolve almost immediately to create race with Then()
			asyncResolve("concurrent", nil)

			select {
			case val := <-done:
				assert.Equal(t, "concurrent!", val, "iteration %d: expected 'concurrent!'", i)
			case <-time.After(3 * time.Second):
				t.Fatalf("timeout on iteration %d", i)
			}
		}
	})

	t.Run("multiple concurrent Then on settling promise", func(t *testing.T) {
		var asyncResolve func(interface{}, error)
		p := NewPromise(func(resolve, reject func(interface{}, error)) {
			asyncResolve = resolve
		})

		const n = 50
		done := make(chan string, n)

		for i := 0; i < n; i++ {
			p.Then(func(v interface{}) (interface{}, error) {
				done <- v.(string)
				return v, nil
			}, nil)
		}

		// Resolve to trigger all subscribers
		asyncResolve("broadcast", nil)

		// Collect all results via channel
		for i := 0; i < n; i++ {
			select {
			case val := <-done:
				assert.Equal(t, "broadcast", val, "subscriber %d should receive 'broadcast'", i)
			case <-time.After(3 * time.Second):
				t.Fatalf("timeout waiting for Then subscriber %d", i)
			}
		}
	})
}
