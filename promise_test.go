package vowlink

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromise_Then(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return value.(string) + " vowlink", nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink", result.value, "Expected value to be 'Hello, World! vowlink'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Then(nil, func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected reason to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onFulfilled and onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(nil, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Then Chain", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return value.(string) + " vowlink", nil
		}, nil).Then(func(value any) (any, error) {
			return value.(string) + "!", nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink!", result.value, "Expected value to be 'Hello, World! vowlink!'")
	})

	// 当.then中返回的不是promise对象时（包括undefined），p2的状态一直都是fulfilled，且值为undefined
	t.Run("Then Chain with Rejection", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return value.(string) + " vowlink", nil
		}, nil).Then(func(value any) (any, error) {
			return value.(string) + "!", nil
		}, func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink!", result.value, "Expected value to be 'Hello, World! vowlink!'")
	})

	t.Run("Then return a Promise with resolve", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return NewPromise(func(resolve func(any, error), reject func(any, error)) {
				resolve(value.(string)+" vowlink", nil)
			}), nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World! vowlink", result.value.(*Promise).GetValue(), "Expected value to be 'Hello, World! vowlink'")
	})

	t.Run("Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return NewPromise(func(resolve func(any, error), reject func(any, error)) {
				reject(nil, errors.New("Something went wrong"))
			}), nil
		}, nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, Rejected, result.value.(*Promise).state, "Expected value to be Rejected")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).GetReason().Error(), "Expected reason to be 'Something went wrong'")
	})

	t.Run("One Then onRejected after Then return a Promise with reject", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Then(func(value any) (any, error) {
			return NewPromise(func(resolve func(any, error), reject func(any, error)) {
				reject(nil, errors.New("Something went wrong"))
			}), nil
		}, nil).Then(nil, func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.value.(*Promise).reason.Error(), "Expected reason to be 'Something went wrong', Then(nil, func(reason error) error) not work")
	})
}

func TestPromise_Catch(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})

	t.Run("Rejected state", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Handled error: Something went wrong", result.reason.Error(), "Expected value to be 'Handled error: Something went wrong'")
	})

	t.Run("Nil onRejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		})

		result := p.Catch(nil)

		assert.Equal(t, Rejected, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Something went wrong", result.reason.Error(), "Expected value to be 'Something went wrong'")
	})
}

func TestPromise_Finally(t *testing.T) {
	t.Run("Fulfilled state", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		})

		result := p.Finally(nil)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Hello, World!", result.value, "Expected value to be 'Hello, World!'")
	})
}

func TestMethod_All(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []any{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 2 rejected", result.reason.Error(), "Expected reason to be 'Promise 2 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := All(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_Any(t *testing.T) {
	t.Run("Any promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := Any(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, "Promise 1", result.value, "Expected value to be 'Promise 1'")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected value to be 'Promise 1 rejected'")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := Race(p1, p2, p3)

		assert.Equal(t, Rejected, result.state, "Expected state to be Rejected")
		assert.Equal(t, "Promise 1 rejected", result.reason.Error(), "Expected reason to be 'Promise 1 rejected'")
	})
}

func TestPromise_AllSettled(t *testing.T) {
	t.Run("All promises fulfilled", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 2", nil)
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []any{"Promise 1", "Promise 2", "Promise 3"}, result.value, "Expected value to be ['Promise 1', 'Promise 2', 'Promise 3']")
	})

	t.Run("One promise rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 1", nil)
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Promise 3", nil)
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []any{"Promise 1", errors.New("Promise 2 rejected"), "Promise 3"}, result.value, "Expected value to be ['Promise 1', errors.New('Promise 2 rejected'), 'Promise 3']")
	})

	t.Run("All promises rejected", func(t *testing.T) {
		p1 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 1 rejected"))
		})

		p2 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 2 rejected"))
		})

		p3 := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Promise 3 rejected"))
		})

		result := AllSettled(p1, p2, p3)

		assert.Equal(t, Fulfilled, result.state, "Expected state to be Fulfilled")
		assert.Equal(t, []any{errors.New("Promise 1 rejected"), errors.New("Promise 2 rejected"), errors.New("Promise 3 rejected")}, result.value, "Expected value to be [errors.New('Promise 1 rejected'), errors.New('Promise 2 rejected'), errors.New('Promise 3 rejected')]")
	})
}

func TestPromise_MultiCatch(t *testing.T) {
	t.Run("Rejected Multi Catch with New Error", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 3 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong", p.GetReason().Error(), "Expected reason to be 'Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Return Value", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return "Recovered value", nil
		}).Then(func(data any) (any, error) {
			return data, nil
		}, nil)

		assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
		assert.Nil(t, p.GetReason(), "Expected reason to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Then return error", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return "Recovered value", nil
		}).Then(func(data any) (any, error) {
			return nil, errors.New("Then error: " + data.(string))
		}, nil).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 2 error: Then error: Recovered value", p.GetReason().Error(), "Expected reason to be 'Handled 2 error: Then error: Recovered value'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with New Error and Finally", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 2 error: " + reason.Error())
		}).Finally(func() error {
			fmt.Println("Finally called")
			return nil
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 3 error: " + reason.Error())
		})

		assert.Equal(t, "Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong", p.GetReason().Error(), "Expected reason to be 'Handled 3 error: Handled 2 error: Handled 1 error: Something went wrong'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")
	})

	t.Run("Rejected Multi Catch with Recover and Finally", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Catch(func(reason error) (any, error) {
			return nil, errors.New("Handled 1 error: " + reason.Error())
		}).Catch(func(reason error) (any, error) {
			return "Recovered value", nil
		}).Finally(func() error {
			fmt.Println("Finally called")
			return nil
		}).Then(func(data any) (any, error) {
			return data, nil
		}, nil)

		assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
		assert.Nil(t, p.GetReason(), "Expected reason to be nil")
	})
}

func TestPromise_ResolveWithError(t *testing.T) {
	p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
		resolve(nil, errors.New("Something went wrong"))
	}).Catch(func(reason error) (any, error) {
		return nil, errors.New("Handled error: " + reason.Error())
	}).Catch(func(reason error) (any, error) {
		return "Recovered value", nil
	}).Finally(func() error {
		fmt.Println("Finally called")
		return nil
	}).Then(func(data any) (any, error) {
		return data, nil
	}, nil)

	assert.Equal(t, "Recovered value", p.GetValue(), "Expected value to be 'Recovered value'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
}

func TestPromise_ResolveWithErrorData(t *testing.T) {
	p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
		resolve(errors.New("Something went wrong"), nil)
	}).Then(func(data any) (any, error) {
		return data.(error).Error(), nil
	}, func(error) (any, error) {
		return nil, errors.New("Handled error")
	}).Catch(func(reason error) (any, error) {
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil
	})

	assert.Equal(t, "Something went wrong", p.GetValue().(string), "Expected value to be 'Something went wrong'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
}

func TestPromise_RejectWithNil(t *testing.T) {
	p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
		reject("Something went wrong", nil)
	}).Then(func(data any) (any, error) {
		return data, nil
	}, func(error) (any, error) {
		return nil, errors.New("Handled error")
	}).Catch(func(reason error) (any, error) {
		return fmt.Sprintf("Recovered value: %v", reason.Error()), nil
	})

	assert.Equal(t, "Recovered value: Handled error", p.GetValue().(string), "Expected value to be 'Recovered value: Handled error'")
	assert.Nil(t, p.GetReason(), "Expected reason to be nil")
}

func TestPromise_FinallyWithError(t *testing.T) {
	t.Run("Finally with error and resolved", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("Hello, World!", nil)
		}).Finally(func() error {
			return errors.New("Finally error")
		}).Then(func(data any) (any, error) {
			return data.(string) + " vowlink", nil
		}, func(reason error) (any, error) {
			return nil, errors.New("Handled error: " + reason.Error())
		})

		assert.Equal(t, "Handled error: Finally error", p.GetReason().Error(), "Expected reason to be 'Handled error: Finally error'")
		assert.Nil(t, p.GetValue(), "Expected value to be nil")

	})

	t.Run("Finally with error and rejected", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			reject(nil, errors.New("Something went wrong"))
		}).Finally(func() error {
			return errors.New("Finally error")
		}).Then(func(data any) (any, error) {
			return data.(string) + " vowlink", nil
		}, func(reason error) (any, error) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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

		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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

			_ = NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve("first", nil)
			resolve("second", nil)           // should not change state
			reject(nil, errors.New("error")) // should not change state
		})

		assert.Equal(t, Fulfilled, p.state)
		assert.Equal(t, "first", p.value)
		assert.Nil(t, p.reason)
	})

	t.Run("rejected state cannot be changed", func(t *testing.T) {
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
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
		assert.Equal(t, []any{}, result.value)
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
		assert.Equal(t, []any{}, result.value)
	})
}

func TestPromise_AsyncThen(t *testing.T) {
	t.Run("async resolve with Then chain", func(t *testing.T) {
		done := make(chan struct{})
		var result string
		p := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				resolve("async", nil)
			}()
		})
		p.Then(func(value any) (any, error) {
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
		p := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				reject(nil, errors.New("async error"))
			}()
		})
		p.Catch(func(reason error) (any, error) {
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
		p := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve(42, nil)
			}()
		})
		p.Then(func(value any) (any, error) {
			r1 = value.(int) * 2
			close(done1)
			return nil, nil
		}, nil)
		p.Then(func(value any) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("first", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve("second", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("third", nil)
			}()
		})
		result := All(p1, p2, p3)
		result.Then(func(value any) (any, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for All async")
		}
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]any)
		assert.Equal(t, 3, len(values))
		assert.Equal(t, "first", values[0])
		assert.Equal(t, "second", values[1])
		assert.Equal(t, "third", values[2])
	})

	t.Run("one async reject fails All", func(t *testing.T) {
		done := make(chan struct{})
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("ok", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("fast fail"))
			}()
		})
		result := All(p1, p2)
		result.Then(nil, func(reason error) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				resolve("slow", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				resolve("fast", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(300 * time.Millisecond)
				resolve("slowest", nil)
			}()
		})
		result := Race(p1, p2, p3)
		result.Then(func(value any) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				resolve("slow", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				reject(nil, errors.New("fast error"))
			}()
		})
		result := Race(p1, p2)
		result.Then(nil, func(reason error) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				reject(nil, errors.New("err1"))
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				resolve("winner", nil)
			}()
		})
		p3 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("loser", nil)
			}()
		})
		result := Any(p1, p2, p3)
		result.Then(func(value any) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("async err1"))
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				reject(nil, errors.New("async err2"))
			}()
		})
		result := Any(p1, p2)
		result.Then(nil, func(reason error) (any, error) {
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
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				resolve("a", nil)
			}()
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				reject(nil, errors.New("err b"))
			}()
		})
		p3 := NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(150 * time.Millisecond)
				resolve("c", nil)
			}()
		})
		result := AllSettled(p1, p2, p3)
		result.Then(func(value any) (any, error) {
			close(done)
			return nil, nil
		}, nil)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("timeout waiting for AllSettled")
		}
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]any)
		assert.Equal(t, 3, len(values))
		assert.Equal(t, "a", values[0])
		assert.EqualError(t, errors.New("err b"), values[1].(error).Error())
		assert.Equal(t, "c", values[2])
	})
}

func TestPromise_NilPromise(t *testing.T) {
	t.Run("All with nil in middle", func(t *testing.T) {
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			resolve("a", nil)
		})
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			resolve("b", nil)
		})
		result := All(p1, nil, p2)
		assert.Equal(t, Fulfilled, result.state)
		values := result.value.([]any)
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "a", values[0])
		assert.Equal(t, "b", values[1])
	})

	t.Run("All with only nil", func(t *testing.T) {
		result := All(nil)
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []any{}, result.value)
	})

	t.Run("AllSettled with nil nil", func(t *testing.T) {
		result := AllSettled(nil, nil)
		assert.Equal(t, Fulfilled, result.state)
		assert.Equal(t, []any{}, result.value)
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
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				go func() {
					defer wg.Done()
					time.Sleep(time.Duration(i%10) * time.Millisecond)
					resolve(i, nil)
				}()
			})
		}
		result := All(promises...)
		done := make(chan struct{})
		result.Then(func(value any) (any, error) {
			close(done)
			return nil, nil
		}, func(reason error) (any, error) {
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
		values := result.value.([]any)
		assert.Equal(t, n, len(values))
		for i := 0; i < n; i++ {
			assert.Equal(t, i, values[i].(int))
		}
	})

	t.Run("concurrent All with mixed timing", func(t *testing.T) {
		const n = 200
		promises := make([]*Promise, n)
		for i := 0; i < n; i++ {
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				go func() {
					time.Sleep(time.Duration(i%5) * time.Millisecond)
					if i%7 == 0 {
						reject(nil, errors.New("err"))
					} else {
						resolve(i, nil)
					}
				}()
			})
		}
		result := All(promises...)
		done := make(chan struct{})
		result.Then(func(value any) (any, error) {
			close(done)
			return nil, nil
		}, func(reason error) (any, error) {
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
		p := NewPromise(func(resolve, reject func(any, error)) {
			panic("executor boom")
		})
		assert.Equal(t, Rejected, p.state)
		assert.NotNil(t, p.GetReason())
		assert.Contains(t, p.GetReason().Error(), "promise executor panic")
		assert.Contains(t, p.GetReason().Error(), "executor boom")
	})

	t.Run("executor panic does not block Then chain", func(t *testing.T) {
		p := NewPromise(func(resolve, reject func(any, error)) {
			panic("executor boom")
		}).Catch(func(reason error) (any, error) {
			return "recovered from: " + reason.Error(), nil
		})
		assert.Equal(t, Fulfilled, p.state)
		assert.Contains(t, p.GetValue().(string), "recovered from")
	})
}

func TestPromise_SubscriberPanicProtection(t *testing.T) {
	t.Run("panicking subscriber does not block others", func(t *testing.T) {
		var asyncResolve func(any, error)
		p := NewPromise(func(resolve, reject func(any, error)) {
			asyncResolve = resolve
		})

		// First subscriber will panic
		result1 := p.Then(func(value any) (any, error) {
			panic("subscriber 1 panic")
		}, nil)

		// Second subscriber should still be called
		done := make(chan struct{})
		var result2Value any
		_ = p.Then(func(value any) (any, error) {
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
		var asyncResolve func(any, error)
		p := NewPromise(func(resolve, reject func(any, error)) {
			asyncResolve = resolve
		})

		// First subscriber panics
		result1 := p.Then(func(value any) (any, error) {
			panic("subscriber 1 boom")
		}, nil)

		// Second subscriber panics
		result2 := p.Then(func(value any) (any, error) {
			panic("subscriber 2 boom")
		}, nil)

		// Third subscriber succeeds
		done := make(chan struct{})
		result3 := p.Then(func(value any) (any, error) {
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
		var asyncReject func(any, error)
		p := NewPromise(func(resolve, reject func(any, error)) {
			asyncReject = reject
		})

		// First subscriber panics
		result1 := p.Then(nil, func(reason error) (any, error) {
			panic("reject handler panic")
		})

		// Second subscriber handles rejection normally
		done := make(chan struct{})
		result2 := p.Then(nil, func(reason error) (any, error) {
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
		var asyncResolve func(any, error)
		source := NewPromise(func(resolve, reject func(any, error)) {
			asyncResolve = resolve
		})

		// 用 channel 追踪每个 subscriber 是否被调用
		const n = 5
		called := make([]chan struct{}, n)
		for i := 0; i < n; i++ {
			called[i] = make(chan struct{})
		}

		// 5 个下游 Promise — idx 在各函数字面量内独立声明，闭包各自绑定独立索引；
		// Go 1.22+ 循环变量已按迭代隔离，即使改写为循环也无需手动捕获
		downstream := make([]*Promise, n)

		// subscriber #1: panic
		func() {
			idx := 0
			downstream[idx] = source.Then(func(value any) (any, error) {
				close(called[idx])
				panic("boom from subscriber 1")
			}, nil)
		}()

		// subscriber #2: 正常
		func() {
			idx := 1
			downstream[idx] = source.Then(func(value any) (any, error) {
				close(called[idx])
				return fmt.Sprintf("sub2 got %v", value), nil
			}, nil)
		}()

		// subscriber #3: panic
		func() {
			idx := 2
			downstream[idx] = source.Then(func(value any) (any, error) {
				close(called[idx])
				panic("boom from subscriber 3")
			}, nil)
		}()

		// subscriber #4: 正常
		func() {
			idx := 3
			downstream[idx] = source.Then(func(value any) (any, error) {
				close(called[idx])
				return fmt.Sprintf("sub4 got %v", value), nil
			}, nil)
		}()

		// subscriber #5: panic
		func() {
			idx := 4
			downstream[idx] = source.Then(func(value any) (any, error) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve(1, nil)
		})

		// Act: 5000 层同步 Then 链，每层将值 +1
		const layers = 5000
		assert.NotPanics(t, func() {
			for i := 0; i < layers; i++ {
				p = p.Then(func(value any) (any, error) {
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
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			resolve(1, nil)
		})

		const rejectAt = 2500
		successCallCount := 0

		assert.NotPanics(t, func() {
			for i := 0; i < 5000; i++ {
				p = p.Then(func(value any) (any, error) {
					successCallCount++
					if i == rejectAt {
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

func TestPromise_DeepPendingChainSettle(t *testing.T) {
	t.Run("10000 layer Then chain built on pending root", func(t *testing.T) {
		var asyncResolve func(any, error)
		p := NewPromise(func(resolve func(any, error), reject func(any, error)) {
			asyncResolve = resolve
		})

		const layers = 10000
		assert.NotPanics(t, func() {
			for i := 0; i < layers; i++ {
				p = p.Then(func(value any) (any, error) {
					return value.(int) + 1, nil
				}, nil)
			}
		}, "building 10000-layer Then chain on pending root should not panic")

		assert.Equal(t, Pending, p.getState(),
			"chain tail should still be Pending before root settles")

		assert.NotPanics(t, func() {
			asyncResolve(0, nil)
		}, "settle cascade through 10000-layer pending chain should not panic")

		assert.Equal(t, Fulfilled, p.getState(),
			"final Promise should be Fulfilled after cascade settle")
		assert.Equal(t, layers, p.GetValue(),
			"final value should be 10000 (0 initial + 10000 increments)")
		assert.Nil(t, p.GetReason(),
			"final Promise should have no rejection reason")
	})
}

// TestPromise_NilPromiseChainCalls
// 场景: NewPromise(nil) 返回一个 Rejected Promise (Round 1 已修复)，
//
//	验证其 Then/Catch/Finally 链式调用正常工作。
//
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
				func(value any) (any, error) {
					t.Fatal("success handler should not be called on Rejected Promise")
					return nil, nil
				},
				func(reason error) (any, error) {
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
			result := nilPromise.Catch(func(reason error) (any, error) {
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
				Catch(func(reason error) (any, error) {
					// 从 nil handler 错误中恢复
					return "recovered from nil handler", nil
				}).
				Finally(func() error {
					finallyCalled = true
					return nil
				}).
				Then(func(value any) (any, error) {
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
				Catch(func(reason error) (any, error) {
					return nil, errors.New("new error from Catch")
				}).
				Finally(func() error {
					return nil
				}).
				Then(func(value any) (any, error) {
					t.Fatal("success handler should not be called on error propagation path")
					return nil, nil
				}, func(reason error) (any, error) {
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
			var asyncResolve func(any, error)
			p := NewPromise(func(resolve, reject func(any, error)) {
				asyncResolve = resolve
			})

			// Use channel to signal completion from Then callback
			done := make(chan string, 1)
			p.Then(func(v any) (any, error) {
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
		var asyncResolve func(any, error)
		p := NewPromise(func(resolve, reject func(any, error)) {
			asyncResolve = resolve
		})

		const n = 50
		done := make(chan string, n)

		for i := 0; i < n; i++ {
			p.Then(func(v any) (any, error) {
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

// TestPromise_RetainedResolveCannotAffectOtherPromises 回归测试（D1）：
// handler 保留 resolve 引用并在 Promise 同步结算后再次调用属于合法用法，
// 延迟调用必须是对原 Promise 的无操作（重复结算被忽略），
// 绝不能决议任何其他 Promise。
func TestPromise_RetainedResolveCannotAffectOtherPromises(t *testing.T) {
	const rounds = 20
	for round := 0; round < rounds; round++ {
		var savedResolve func(any, error)
		p1 := NewPromise(func(resolve, reject func(any, error)) {
			savedResolve = resolve
			resolve("p1-value", nil)
		})
		assert.Equal(t, Fulfilled, p1.getState(), "round %d: p1 should be Fulfilled", round)
		assert.Equal(t, "p1-value", p1.GetValue(), "round %d: p1 initial value", round)

		// 创建并同步结算一个无关 Promise，再创建一个保持 Pending 的无关 Promise
		p2 := NewPromise(func(resolve, reject func(any, error)) {
			resolve("p2-value", nil)
		})
		p3 := NewPromise(func(resolve, reject func(any, error)) {
			// 刻意保持 Pending
		})

		// 延迟调用：按 JS Promise 语义必须是针对已结算 p1 的无操作
		savedResolve("late-call", nil)

		assert.Equal(t, "p1-value", p1.GetValue(), "round %d: p1 value must be unchanged by late resolve", round)
		assert.Nil(t, p1.GetReason(), "round %d: p1 reason must remain nil", round)
		assert.Equal(t, Fulfilled, p2.getState(), "round %d: settled p2 state must not change", round)
		assert.Equal(t, "p2-value", p2.GetValue(), "round %d: settled p2 value must not change", round)
		assert.Equal(t, Pending, p3.getState(), "round %d: pending p3 must not be settled by late resolve", round)
		assert.Nil(t, p3.GetValue(), "round %d: pending p3 value must remain nil", round)
		assert.Nil(t, p3.GetReason(), "round %d: pending p3 reason must remain nil", round)
	}
}

// TestAggregateError_ConcurrentError 回归测试（D2）：
// 对填充后的 *AggregateError 并发调用 Error() 与 InvalidateError()
// 不得产生数据竞争（由 -race 验证），且 Error() 返回内容始终正确。
func TestAggregateError_ConcurrentError(t *testing.T) {
	ae := NewAggregateError(4)
	for i := 0; i < 4; i++ {
		ae.Errors = append(ae.Errors, fmt.Errorf("err-%d", i))
	}
	const expected = "All promises were rejected: err-0, err-1, err-2, err-3"

	const goroutines = 8
	const iterations = 2000
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if got := ae.Error(); got != expected {
					t.Errorf("concurrent Error() returned %q, want %q", got, expected)
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ae.InvalidateError()
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, expected, ae.Error())
}

// TestPromise_AnyErrorsInInputOrder 回归测试（D6）：
// 输入按 0,1,2,3,4 顺序，但以逆序/乱序错峰完成（完成顺序 1,3,4,2,0），
// AggregateError.Errors 必须严格按输入顺序排列。
func TestPromise_AnyErrorsInInputOrder(t *testing.T) {
	delays := []time.Duration{100, 20, 80, 40, 60}
	promises := make([]*Promise, len(delays))
	for i := range delays {
		promises[i] = NewPromise(func(resolve, reject func(any, error)) {
			go func() {
				time.Sleep(delays[i])
				reject(nil, fmt.Errorf("err-%d", i))
			}()
		})
	}

	done := make(chan struct{})
	result := Any(promises...)
	result.Then(nil, func(reason error) (any, error) {
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
	assert.True(t, ok, "Expected reason to be an *AggregateError")
	assert.Equal(t, len(delays), len(aggErr.Errors))
	assert.Equal(t, "err-0", aggErr.Errors[0].Error())
	assert.Equal(t, "err-1", aggErr.Errors[1].Error())
	assert.Equal(t, "err-2", aggErr.Errors[2].Error())
	assert.Equal(t, "err-3", aggErr.Errors[3].Error())
	assert.Equal(t, "err-4", aggErr.Errors[4].Error())
}

// settleFanIn 让 workers 个 goroutine 在同一时刻开始结算 settlers 中的 Promise，
// 模拟大规模并发扇入结算。
func settleFanIn(settlers []func(), workers int) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			<-start
			for j := base; j < len(settlers); j += workers {
				settlers[j]()
			}
		}(w)
	}
	close(start)
	wg.Wait()
}

// TestPromise_RaceConcurrentFanIn 回归测试：
// 16 个 goroutine 同时结算 128 个混合 fulfill/reject 的输入，
// 重复 50 轮（配合 -race），Race 结果必须属于输入集合。
func TestPromise_RaceConcurrentFanIn(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 50
	)

	allowedValues := make(map[int]struct{}, n)
	allowedReasons := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		if i%3 == 0 {
			allowedReasons[fmt.Sprintf("fanin err %d", i)] = struct{}{}
		} else {
			allowedValues[i] = struct{}{}
		}
	}

	for round := 0; round < rounds; round++ {
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		for i := 0; i < n; i++ {
			if i%3 == 0 {
				reason := fmt.Errorf("fanin err %d", i)
				promises[i] = NewPromise(func(resolve, reject func(any, error)) {
					settlers[i] = func() { reject(nil, reason) }
				})
			} else {
				promises[i] = NewPromise(func(resolve, reject func(any, error)) {
					settlers[i] = func() { resolve(i, nil) }
				})
			}
		}

		result := Race(promises...)
		settled := make(chan struct{})
		result.Then(func(value any) (any, error) {
			close(settled)
			return nil, nil
		}, func(reason error) (any, error) {
			close(settled)
			return nil, nil
		})

		settleFanIn(settlers, workers)

		select {
		case <-settled:
		case <-time.After(3 * time.Second):
			t.Fatalf("round %d: timeout waiting for Race fan-in settle", round)
		}

		switch result.getState() {
		case Fulfilled:
			v, ok := result.value.(int)
			if !assert.True(t, ok, "round %d: Race value must be int", round) {
				continue
			}
			_, allowed := allowedValues[v]
			assert.True(t, allowed, "round %d: Race value %d not in input set", round, v)
		case Rejected:
			_, allowed := allowedReasons[result.reason.Error()]
			assert.True(t, allowed, "round %d: Race reason %q not in input set", round, result.reason.Error())
		default:
			t.Fatalf("round %d: Race still pending after all inputs settled", round)
		}
	}
}

// TestPromise_AllConcurrentFanIn 回归测试：
// 16 个 goroutine 同时结算 128 个输入，重复 50 轮（配合 -race）。
// 全部 fulfill 时结果必须是完整有序值；混合 reject 时结果必须是首个拒绝错误之一。
func TestPromise_AllConcurrentFanIn(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 50
	)

	t.Run("all fulfilled values in order", func(t *testing.T) {
		for round := 0; round < rounds; round++ {
			promises := make([]*Promise, n)
			settlers := make([]func(), n)
			for i := 0; i < n; i++ {
				promises[i] = NewPromise(func(resolve, reject func(any, error)) {
					settlers[i] = func() { resolve(i, nil) }
				})
			}

			result := All(promises...)
			settled := make(chan struct{})
			result.Then(func(value any) (any, error) {
				close(settled)
				return nil, nil
			}, func(reason error) (any, error) {
				close(settled)
				return nil, nil
			})

			settleFanIn(settlers, workers)

			select {
			case <-settled:
			case <-time.After(3 * time.Second):
				t.Fatalf("round %d: timeout waiting for All fan-in settle", round)
			}

			assert.Equal(t, Fulfilled, result.getState(), "round %d", round)
			values, ok := result.value.([]any)
			if !assert.True(t, ok, "round %d: All value must be []any", round) {
				continue
			}
			assert.Equal(t, n, len(values), "round %d", round)
			for i := 0; i < n; i++ {
				assert.Equal(t, i, values[i].(int), "round %d index %d", round, i)
			}
		}
	})

	t.Run("mixed rejects settled by first rejection", func(t *testing.T) {
		allowedReasons := make(map[string]struct{})
		for i := 0; i < n; i++ {
			if i%5 == 0 {
				allowedReasons[fmt.Sprintf("all err %d", i)] = struct{}{}
			}
		}

		for round := 0; round < rounds; round++ {
			promises := make([]*Promise, n)
			settlers := make([]func(), n)
			for i := 0; i < n; i++ {
				if i%5 == 0 {
					reason := fmt.Errorf("all err %d", i)
					promises[i] = NewPromise(func(resolve, reject func(any, error)) {
						settlers[i] = func() { reject(nil, reason) }
					})
				} else {
					promises[i] = NewPromise(func(resolve, reject func(any, error)) {
						settlers[i] = func() { resolve(i, nil) }
					})
				}
			}

			result := All(promises...)
			settled := make(chan struct{})
			result.Then(func(value any) (any, error) {
				close(settled)
				return nil, nil
			}, func(reason error) (any, error) {
				close(settled)
				return nil, nil
			})

			settleFanIn(settlers, workers)

			select {
			case <-settled:
			case <-time.After(3 * time.Second):
				t.Fatalf("round %d: timeout waiting for All fan-in settle", round)
			}

			assert.Equal(t, Rejected, result.getState(), "round %d", round)
			_, allowed := allowedReasons[result.reason.Error()]
			assert.True(t, allowed, "round %d: All reason %q must be one of the input rejections", round, result.reason.Error())
		}
	})
}

// customCodeErr 用于 errors.As 互操作测试的自定义错误类型。
type customCodeErr struct{ code int }

func (e *customCodeErr) Error() string { return fmt.Sprintf("custom-%d", e.code) }

// TestPromise_AggregateErrorUnwrapInterop 回归测试（P1-1）：
// *AggregateError 实现 Unwrap() []error（Go 1.20+ 多错误惯例），
// errors.Is/errors.As 必须能遍历并命中成员错误中的 sentinel 与自定义类型，
// 覆盖直接 NewAggregateError 构造与 Any 全拒 reason 两条路径。
func TestPromise_AggregateErrorUnwrapInterop(t *testing.T) {
	sentinel := errors.New("unwrap-sentinel")
	custom := &customCodeErr{code: 42}

	t.Run("direct NewAggregateError construction", func(t *testing.T) {
		aggErr := NewAggregateError(3)
		aggErr.Errors = append(aggErr.Errors,
			errors.New("plain"),
			fmt.Errorf("ctx: %w", sentinel),
			custom,
		)

		if !errors.Is(aggErr, sentinel) {
			t.Errorf("errors.Is(aggErr, sentinel)=false，遍历未命中成员中 wrap 的 sentinel")
		}
		var target *customCodeErr
		if !errors.As(aggErr, &target) {
			t.Errorf("errors.As(aggErr, **customCodeErr)=false，遍历未命中成员中的自定义类型")
		}
		if target != nil && target != custom {
			t.Errorf("errors.As 命中 %v, want 同一成员实例 %v", target, custom)
		}
		if errors.Is(aggErr, errors.New("non-member")) {
			t.Errorf("errors.Is 命中非成员错误，遍历语义异常")
		}
	})

	t.Run("Any all-rejected reason", func(t *testing.T) {
		p1 := NewPromise(func(res, rej func(any, error)) { rej(nil, errors.New("plain")) })
		p2 := NewPromise(func(res, rej func(any, error)) { rej(nil, fmt.Errorf("ctx: %w", sentinel)) })
		p3 := NewPromise(func(res, rej func(any, error)) { rej(nil, custom) })

		result := Any(p1, p2, p3)
		assert.Equal(t, Rejected, result.getState())
		reason := result.GetReason()
		if reason == nil {
			t.Fatal("Any 全拒应返回非 nil reason")
		}

		if !errors.Is(reason, sentinel) {
			t.Errorf("errors.Is(Any 的 reason, sentinel)=false — 成员错误明确 wrap 了 sentinel，标准错误遍历静默假阴性")
		}
		var target *customCodeErr
		if !errors.As(reason, &target) {
			t.Errorf("errors.As(Any 的 reason, **customCodeErr)=false — 无法定位成员中的自定义错误类型")
		}
		if target != nil && target != custom {
			t.Errorf("errors.As 命中 %v, want 同一成员实例 %v", target, custom)
		}
	})
}

// fanInSubset 让 workers 个 goroutine 在同一时刻开始结算 settlers 中 idxs 指定位置的 Promise，
// 用于混合窗口 fan-in 场景（部分输入在组合器构建期间已并发结算，其余构建后真并发扇入）。
func fanInSubset(settlers []func(), idxs []int, workers int) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			<-start
			for j := base; j < len(idxs); j += workers {
				settlers[idxs[j]]()
			}
		}(w)
	}
	close(start)
	wg.Wait()
}

// rangeInts 返回 [from, to) 区间的整数序列。
func rangeInts(from, to int) []int {
	out := make([]int, 0, to-from)
	for i := from; i < to; i++ {
		out = append(out, i)
	}
	return out
}

// TestPromise_ConcurrentRetainedResolve 回归测试（D1 并发扩展）：
// 8 个 goroutine 各创建 150 个"handler 保留 resolve 引用并同步结算"的 Promise，
// 交错创建无关噪声 Promise（已结算 + 永久 pending）模拟池复用最坏时序，
// 最后全量延迟调用保留的 resolve（毒化值）。
// 任何跨 Promise 决议（废池前的 D1 缺陷模式）都会污染断言。
func TestPromise_ConcurrentRetainedResolve(t *testing.T) {
	const (
		goroutines = 8
		perG       = 150
	)
	type retained struct {
		p    *Promise
		want string
		call func(any, error)
	}
	all := make([][]retained, goroutines)
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			items := make([]retained, 0, perG)
			rng := rand.New(rand.NewSource(int64(g)))
			for i := 0; i < perG; i++ {
				want := fmt.Sprintf("g%d-i%d", g, i)
				var saved func(any, error)
				p := NewPromise(func(resolve, reject func(any, error)) {
					saved = resolve
					resolve(want, nil)
				})
				items = append(items, retained{p: p, want: want, call: saved})
				// 交错创建无关 Promise（已结算 + 永久 pending），模拟池复用时的最坏时序
				_ = NewPromise(func(resolve, reject func(any, error)) {
					resolve(fmt.Sprintf("noise-%d", rng.Int()), nil)
				})
				if i%3 == 0 {
					_ = NewPromise(func(resolve, reject func(any, error)) {})
				}
			}
			all[g] = items
		}(g)
	}
	wg.Wait()

	// 延迟调用全部保留的 resolve（毒化值），任何跨 Promise 决议都会污染断言
	for g := range all {
		for i := range all[g] {
			all[g][i].call("LATE-POISON", nil)
		}
	}
	for g := range all {
		for i := range all[g] {
			it := all[g][i]
			if got := it.p.GetValue(); got != it.want {
				t.Fatalf("D1 回归 FAIL: g%d-i%d value=%v, want %q（跨 Promise 污染）", g, i, got, it.want)
			}
			if it.p.getState() != Fulfilled {
				t.Fatalf("D1 回归 FAIL: g%d-i%d state=%v", g, i, it.p.getState())
			}
			if r := it.p.GetReason(); r != nil {
				t.Fatalf("D1 回归 FAIL: g%d-i%d reason=%v", g, i, r)
			}
		}
	}
}

// TestPromise_AnyAllRejectMixedWindowFanIn 回归测试（D6 扩展 + 混合窗口盲区）：
// 一半输入在 Any() 构建期间并发结算（executor 循环中 sync-dispatch 与
// append-subscriber 两条路径随机命中），另一半构建后 16-worker 真 fan-in，25 轮 × 128 输入。
// 断言 AggregateError.Errors 严格按输入序（证明索引写 + 恰好一次分派 + 可见性）。
func TestPromise_AnyAllRejectMixedWindowFanIn(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 25
	)
	for round := 0; round < rounds; round++ {
		rng := rand.New(rand.NewSource(int64(round) + 1))
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		for i := 0; i < n; i++ {
			reason := fmt.Errorf("any err %d", i)
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				settlers[i] = func() { reject(nil, reason) }
			})
		}

		perm := rng.Perm(n)
		during := perm[:n/2]
		after := perm[n/2:]
		durations := make([]time.Duration, len(during))
		for k := range durations {
			durations[k] = time.Duration(rng.Intn(300)) * time.Microsecond
		}

		var bg sync.WaitGroup
		for k, idx := range during {
			bg.Add(1)
			go func(idx int, d time.Duration) {
				defer bg.Done()
				time.Sleep(d)
				settlers[idx]()
			}(idx, durations[k])
		}

		result := Any(promises...)
		settled := make(chan struct{})
		result.Then(nil, func(reason error) (any, error) {
			close(settled)
			return nil, nil
		})

		fanInSubset(settlers, after, workers)
		bg.Wait()

		select {
		case <-settled:
		case <-time.After(5 * time.Second):
			t.Fatalf("round %d: timeout waiting for Any all-reject", round)
		}

		if st := result.getState(); st != Rejected {
			t.Fatalf("round %d: Any state=%v, want Rejected", round, st)
		}
		aggErr, ok := result.GetReason().(*AggregateError)
		if !ok {
			t.Fatalf("round %d: reason type %T, want *AggregateError", round, result.GetReason())
		}
		if len(aggErr.Errors) != n {
			t.Fatalf("round %d: len(Errors)=%d, want %d", round, len(aggErr.Errors), n)
		}
		for i := 0; i < n; i++ {
			want := fmt.Sprintf("any err %d", i)
			if aggErr.Errors[i] == nil || aggErr.Errors[i].Error() != want {
				t.Fatalf("round %d: Errors[%d]=%v, want %q（索引序/恰好一次分派被破坏）", round, i, aggErr.Errors[i], want)
			}
		}
	}
}

// TestPromise_AnyMixedFanInNeverRejects 回归测试（混合窗口盲区）：
// 128 个输入中恰有一个 fulfiller（位置每轮随机），混合窗口并发结算，25 轮。
// 存在 fulfiller 时 Any 绝不允许 reject（早决短路后迟到 onRejected 写 Errors 只写无读）。
func TestPromise_AnyMixedFanInNeverRejects(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 25
	)
	for round := 0; round < rounds; round++ {
		rng := rand.New(rand.NewSource(int64(round) + 1000))
		winner := rng.Intn(n)
		want := fmt.Sprintf("win-%d-%d", round, winner)
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		for i := 0; i < n; i++ {
			if i == winner {
				promises[i] = NewPromise(func(resolve, reject func(any, error)) {
					settlers[i] = func() { resolve(want, nil) }
				})
				continue
			}
			reason := fmt.Errorf("mix err %d", i)
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				settlers[i] = func() { reject(nil, reason) }
			})
		}

		perm := rng.Perm(n)
		during := perm[:n/2]
		after := perm[n/2:]
		durations := make([]time.Duration, len(during))
		for k := range durations {
			durations[k] = time.Duration(rng.Intn(300)) * time.Microsecond
		}

		var bg sync.WaitGroup
		for k, idx := range during {
			bg.Add(1)
			go func(idx int, d time.Duration) {
				defer bg.Done()
				time.Sleep(d)
				settlers[idx]()
			}(idx, durations[k])
		}

		result := Any(promises...)
		settled := make(chan struct{})
		var gotValue any
		var gotReason error
		result.Then(func(v any) (any, error) {
			gotValue = v
			close(settled)
			return nil, nil
		}, func(e error) (any, error) {
			gotReason = e
			close(settled)
			return nil, nil
		})

		fanInSubset(settlers, after, workers)
		bg.Wait()

		select {
		case <-settled:
		case <-time.After(5 * time.Second):
			t.Fatalf("round %d: timeout", round)
		}

		if st := result.getState(); st != Fulfilled {
			t.Fatalf("round %d: Any state=%v reason=%v, want Fulfilled=%q（存在 fulfiller 时绝不允许 reject）", round, st, result.GetReason(), want)
		}
		if gotReason != nil {
			t.Fatalf("round %d: unexpected reason %v", round, gotReason)
		}
		if gotValue != want || result.GetValue() != want {
			t.Fatalf("round %d: value=%v/%v, want %q", round, gotValue, result.GetValue(), want)
		}
	}
}

// TestPromise_AllSettledMixedWindowFanIn 回归测试（C2 盲区补强）：
// AllSettled 首个真并发混合窗口 fan-in（偶数位 resolve、奇数位 reject），25 轮 × 128 输入。
// 结果切片必须严格逐索引等于期望值/错误（索引写 + 原子计数 + 可见性）。
func TestPromise_AllSettledMixedWindowFanIn(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 25
	)
	for round := 0; round < rounds; round++ {
		rng := rand.New(rand.NewSource(int64(round) + 2000))
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		want := make([]any, n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				v := fmt.Sprintf("v-%d-%d", round, i)
				want[i] = v
				promises[i] = NewPromise(func(resolve, reject func(any, error)) {
					settlers[i] = func() { resolve(v, nil) }
				})
				continue
			}
			reason := fmt.Errorf("as err %d-%d", round, i)
			want[i] = reason
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				settlers[i] = func() { reject(nil, reason) }
			})
		}

		perm := rng.Perm(n)
		during := perm[:n/2]
		after := perm[n/2:]
		durations := make([]time.Duration, len(during))
		for k := range durations {
			durations[k] = time.Duration(rng.Intn(300)) * time.Microsecond
		}

		var bg sync.WaitGroup
		for k, idx := range during {
			bg.Add(1)
			go func(idx int, d time.Duration) {
				defer bg.Done()
				time.Sleep(d)
				settlers[idx]()
			}(idx, durations[k])
		}

		result := AllSettled(promises...)
		settled := make(chan struct{})
		result.Then(func(v any) (any, error) {
			close(settled)
			return nil, nil
		}, nil)

		fanInSubset(settlers, after, workers)
		bg.Wait()

		select {
		case <-settled:
		case <-time.After(5 * time.Second):
			t.Fatalf("round %d: timeout", round)
		}
		if st := result.getState(); st != Fulfilled {
			t.Fatalf("round %d: AllSettled state=%v, want Fulfilled", round, st)
		}
		values, ok := result.GetValue().([]any)
		if !ok || len(values) != n {
			t.Fatalf("round %d: value=%v (%T), want []any len %d", round, result.GetValue(), result.GetValue(), n)
		}
		for i := 0; i < n; i++ {
			if values[i] != want[i] {
				t.Fatalf("round %d: values[%d]=%v, want %v（索引写/可见性被破坏）", round, i, values[i], want[i])
			}
		}
	}
}

// TestPromise_AllEarlyRejectLateFulfillSafety 回归测试（短路 + 迟到写安全）：
// p0 先结算使 All 构建时同步短路拒绝，随后 127 个迟到 fulfill 全量 fan-in。
// 迟到回调继续写 values[i] 并递减 pending，pending 永不归零（onRejected 不递减）
// → resolve 永不被调用；All 状态/原因/值必须保持不变，20 轮 × 128 输入。
func TestPromise_AllEarlyRejectLateFulfillSafety(t *testing.T) {
	const (
		n       = 128
		workers = 16
		rounds  = 20
	)
	for round := 0; round < rounds; round++ {
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		firstErr := fmt.Errorf("first err %d", round)
		promises[0] = NewPromise(func(resolve, reject func(any, error)) {
			settlers[0] = func() { reject(nil, firstErr) }
		})
		for i := 1; i < n; i++ {
			v := fmt.Sprintf("late-%d-%d", round, i)
			promises[i] = NewPromise(func(resolve, reject func(any, error)) {
				settlers[i] = func() { resolve(v, nil) }
			})
		}

		// p0 先结算 → All 构建时走 sync-dispatch 立即短路拒绝
		settlers[0]()
		result := All(promises...)
		if st := result.getState(); st != Rejected {
			t.Fatalf("round %d: All 应同步短路拒绝, state=%v", round, st)
		}
		if r := result.GetReason(); r == nil || r.Error() != firstErr.Error() {
			t.Fatalf("round %d: reason=%v, want %v", round, r, firstErr)
		}
		if v := result.GetValue(); v != nil {
			t.Fatalf("round %d: GetValue=%v, want nil after rejection", round, v)
		}

		// 迟到 fulfill 全量 fan-in：迟到回调继续写 values[i] 并递减 pending，
		// pending 永不归零（onRejected 不递减）→ resolve 永不被调用
		fanInSubset(settlers, rangeInts(1, n), workers)

		if st := result.getState(); st != Rejected {
			t.Fatalf("round %d: state flipped to %v after late fulfills", round, st)
		}
		if r := result.GetReason(); r == nil || r.Error() != firstErr.Error() {
			t.Fatalf("round %d: reason changed after late fulfills: %v", round, r)
		}
		if v := result.GetValue(); v != nil {
			t.Fatalf("round %d: GetValue=%v after late fulfills, want nil", round, v)
		}
	}
}

// TestPromise_CombinatorResolveWithError 文档承诺核验（README "Resolve & Reject" 组合器段落）：
// 组合器内 resolve(value, err) 一律按拒绝处理——All/Race 以 err 拒绝且丢弃 value、
// AllSettled 将 err 存入结果切片对应索引、Any 将 err 记入 AggregateError；
// Any 在另一输入 fulfill 时仍必须 Fulfilled。
func TestPromise_CombinatorResolveWithError(t *testing.T) {
	boom := errors.New("boom-payload-err")

	t.Run("All rejects with err and discards value", func(t *testing.T) {
		p := NewPromise(func(res, rej func(any, error)) { res("payload", boom) })
		r := All(p)
		if r.getState() != Rejected {
			t.Fatalf("state=%v, want Rejected", r.getState())
		}
		if r.GetReason() != boom {
			t.Fatalf("reason=%v, want boom (identity)", r.GetReason())
		}
		if r.GetValue() != nil {
			t.Fatalf("value=%v, want nil (payload discarded)", r.GetValue())
		}
	})

	t.Run("Race rejects with err", func(t *testing.T) {
		p := NewPromise(func(res, rej func(any, error)) { res("payload", boom) })
		r := Race(p)
		if r.getState() != Rejected || r.GetReason() != boom {
			t.Fatalf("state=%v reason=%v, want Rejected/boom", r.getState(), r.GetReason())
		}
	})

	t.Run("AllSettled stores err at index", func(t *testing.T) {
		p1 := NewPromise(func(res, rej func(any, error)) { res("payload", boom) })
		p2 := NewPromise(func(res, rej func(any, error)) { res("ok", nil) })
		r := AllSettled(p1, p2)
		if r.getState() != Fulfilled {
			t.Fatalf("state=%v, want Fulfilled", r.getState())
		}
		values, ok := r.GetValue().([]any)
		if !ok || len(values) != 2 {
			t.Fatalf("value=%v", r.GetValue())
		}
		if values[0] != boom {
			t.Fatalf("values[0]=%v, want boom error stored", values[0])
		}
		if values[1] != "ok" {
			t.Fatalf("values[1]=%v, want ok", values[1])
		}
	})

	t.Run("Any records err in AggregateError", func(t *testing.T) {
		p := NewPromise(func(res, rej func(any, error)) { res("payload", boom) })
		r := Any(p)
		if r.getState() != Rejected {
			t.Fatalf("state=%v, want Rejected", r.getState())
		}
		agg, ok := r.GetReason().(*AggregateError)
		if !ok || len(agg.Errors) != 1 || agg.Errors[0] != boom {
			t.Fatalf("reason=%v, want AggregateError[boom]", r.GetReason())
		}
	})

	t.Run("Any still fulfills when another input fulfills", func(t *testing.T) {
		p1 := NewPromise(func(res, rej func(any, error)) { res("payload", boom) })
		p2 := NewPromise(func(res, rej func(any, error)) { res("win", nil) })
		r := Any(p1, p2)
		if r.getState() != Fulfilled || r.GetValue() != "win" {
			t.Fatalf("state=%v value=%v, want Fulfilled/win", r.getState(), r.GetValue())
		}
	})
}

// TestPromise_PanicPathPrefixes 钉住 panic 包装双前缀现状（P2-3，已文档化于 README Concurrency 节）：
// 上游已结算时 Then 走同步快路径，handler panic 由 executor recover 包装为
// "promise executor panic: ..."；上游 pending 时走异步 subscriber 分派，
// 由 dispatchSubscriber recover 包装为 "subscriber callback panic: ..."。
// 两路径行为等价：下游均 Rejected、均含原始 panic 值、均可 Catch 恢复。
func TestPromise_PanicPathPrefixes(t *testing.T) {
	// 同步快路径：上游已结算 → handler panic 由 NewPromise executor recover 兜底
	src1 := NewPromise(func(res, _ func(any, error)) { res("v", nil) })
	d1 := src1.Then(func(any) (any, error) { panic("same-boom") }, nil)

	// 异步 subscriber 路径：上游 pending → handler panic 由 dispatchSubscriber recover 兜底
	var asyncRes func(any, error)
	src2 := NewPromise(func(res, _ func(any, error)) { asyncRes = res })
	d2 := src2.Then(func(any) (any, error) { panic("same-boom") }, nil)
	asyncRes("v", nil)

	for i, d := range []*Promise{d1, d2} {
		if d.getState() != Rejected {
			t.Fatalf("case %d: downstream state=%v, want Rejected", i, d.getState())
		}
		r := d.GetReason()
		if r == nil || !strings.Contains(r.Error(), "same-boom") {
			t.Fatalf("case %d: reason=%v, want contains same-boom", i, r)
		}
	}
	if !strings.Contains(d1.GetReason().Error(), "promise executor panic") {
		t.Errorf("sync 路径前缀与既往审计记录不符: %q", d1.GetReason())
	}
	if !strings.Contains(d2.GetReason().Error(), "subscriber callback panic") {
		t.Errorf("async 路径前缀与既往审计记录不符: %q", d2.GetReason())
	}
	// 两路径下游均可被 Catch 恢复（行为等价性）
	for i, d := range []*Promise{d1, d2} {
		r := d.Catch(func(e error) (any, error) { return "recovered", nil })
		if r.GetValue() != "recovered" {
			t.Fatalf("case %d: Catch recovery failed, value=%v", i, r.GetValue())
		}
	}
}

// TestPromise_ConcurrentReadersDuringFanIn 回归测试（并发读安全）：
// All fan-in 结算期间 4 个 reader goroutine 持续并发读取结果 Promise 与全部输入
// 的 GetValue/GetReason/getState，15 轮 × 64 输入（配合 -race 验证读路径无数据竞争）。
// 存在 reject 输入（i%4==0）→ All 最终必须 Rejected。
func TestPromise_ConcurrentReadersDuringFanIn(t *testing.T) {
	const (
		n       = 64
		workers = 8
		rounds  = 15
		readers = 4
	)
	for round := 0; round < rounds; round++ {
		promises := make([]*Promise, n)
		settlers := make([]func(), n)
		for i := 0; i < n; i++ {
			v := i
			promises[i] = NewPromise(func(res, rej func(any, error)) {
				settlers[i] = func() {
					if v%4 == 0 {
						rej(nil, fmt.Errorf("e%d", v))
					} else {
						res(v, nil)
					}
				}
			})
		}
		result := All(promises...)
		stop := make(chan struct{})
		var rwg sync.WaitGroup
		for r := 0; r < readers; r++ {
			rwg.Add(1)
			go func() {
				defer rwg.Done()
				for {
					select {
					case <-stop:
						return
					default:
						_ = result.GetValue()
						_ = result.GetReason()
						_ = result.getState()
						for _, p := range promises {
							_ = p.GetValue()
							_ = p.getState()
						}
					}
				}
			}()
		}
		settleFanIn(settlers, workers)
		close(stop)
		rwg.Wait()
		// i%4==0 存在 → All 必须 Rejected
		if result.getState() != Rejected {
			t.Fatalf("round %d: state=%v, want Rejected", round, result.getState())
		}
	}
}

// TestPromise_DuplicatePromiseInputs 回归测试：同一 Promise 实例重复作为组合器输入时，
// 每次输入独立订阅、独立占位结果索引（All/AllSettled 值按位重复、Any Errors 按位重复、
// Race 正常选主），不得因指针相同而丢失或串位。
func TestPromise_DuplicatePromiseInputs(t *testing.T) {
	t.Run("All duplicate fulfilled", func(t *testing.T) {
		p := NewPromise(func(res, _ func(any, error)) { res("dup", nil) })
		r := All(p, p, p)
		vals, ok := r.GetValue().([]any)
		if !ok || len(vals) != 3 {
			t.Fatalf("value=%v", r.GetValue())
		}
		for i := range vals {
			if vals[i] != "dup" {
				t.Fatalf("vals[%d]=%v", i, vals[i])
			}
		}
	})
	t.Run("Any duplicate rejected", func(t *testing.T) {
		boom := errors.New("dup-boom")
		p := NewPromise(func(_, rej func(any, error)) { rej(nil, boom) })
		r := Any(p, p)
		agg, ok := r.GetReason().(*AggregateError)
		if !ok || len(agg.Errors) != 2 {
			t.Fatalf("reason=%v", r.GetReason())
		}
		if agg.Errors[0] != boom || agg.Errors[1] != boom {
			t.Fatalf("errors=%v", agg.Errors)
		}
	})
	t.Run("Race duplicate", func(t *testing.T) {
		p := NewPromise(func(res, _ func(any, error)) { res("rd", nil) })
		r := Race(p, p)
		if r.getState() != Fulfilled || r.GetValue() != "rd" {
			t.Fatalf("state=%v value=%v", r.getState(), r.GetValue())
		}
	})
	t.Run("AllSettled duplicate resolve-with-error", func(t *testing.T) {
		boom := errors.New("dup2")
		p := NewPromise(func(res, _ func(any, error)) { res("x", boom) })
		r := AllSettled(p, p)
		vals, ok := r.GetValue().([]any)
		if !ok || len(vals) != 2 {
			t.Fatalf("value=%v", r.GetValue())
		}
		if vals[0] != boom || vals[1] != boom {
			t.Fatalf("vals=%v", vals)
		}
	})
}

// TestPromise_ThenSettleWindowStress 回归测试（Then 注册 × settle 并发窗口）：
// 64 个 goroutine 并发注册 Then 的同时结算源 Promise，压 snapshot→Lock 双检窗口，
// 30 轮。fired 计数必须精确等于注册数（恰好一次分派，无丢失、无双重），
// 所有下游最终必须 Fulfilled 且值正确。
func TestPromise_ThenSettleWindowStress(t *testing.T) {
	const (
		registrars = 64
		rounds     = 30
	)
	for round := 0; round < rounds; round++ {
		var asyncRes func(any, error)
		p := NewPromise(func(res, _ func(any, error)) { asyncRes = res })
		want := fmt.Sprintf("r%d", round)
		done := make(chan struct{}, registrars)
		var fired atomic.Int32
		downstreams := make([]*Promise, registrars)
		var wg sync.WaitGroup
		for k := 0; k < registrars; k++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				downstreams[k] = p.Then(func(v any) (any, error) {
					fired.Add(1)
					done <- struct{}{}
					return v, nil
				}, nil)
			}(k)
		}
		// 与注册并发地结算，压 Then 的 snapshot→Lock 双检窗口
		asyncRes(want, nil)
		for i := 0; i < registrars; i++ {
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatalf("round %d: callback %d never fired（分派丢失）", round, i)
			}
		}
		wg.Wait()
		if got := fired.Load(); got != registrars {
			t.Fatalf("round %d: fired=%d, want %d（恰好一次分派被破坏）", round, got, registrars)
		}
		// 所有下游最终必须 Fulfilled 且值正确
		for k, d := range downstreams {
			deadline := time.Now().Add(2 * time.Second)
			for d.getState() == Pending && time.Now().Before(deadline) {
				runtime.Gosched()
			}
			if d.getState() != Fulfilled || d.GetValue() != want {
				t.Fatalf("round %d: downstream[%d] state=%v value=%v, want Fulfilled/%q", round, k, d.getState(), d.GetValue(), want)
			}
		}
	}
}

// TestPromise_GetValueOnRejectedReturnsPayload 钉住实际行为（P2-1 文档修正配套）：
// 未结算时 GetValue 返回 nil；reject(value, err) 传入的 value 会被保留，
// 即使状态为 Rejected 也可经 GetValue 读取（与 GetValue 文档注释及 README 一致）。
func TestPromise_GetValueOnRejectedReturnsPayload(t *testing.T) {
	pending := NewPromise(func(res, rej func(any, error)) {})
	assert.Nil(t, pending.GetValue(), "未结算 Promise 的 GetValue 应为 nil")

	p := NewPromise(func(res, rej func(any, error)) { rej("payload-on-reject", errors.New("boom")) })
	assert.Equal(t, Rejected, p.getState())
	assert.Equal(t, "payload-on-reject", p.GetValue(), "reject 传入的 value 应被保留并可经 GetValue 读取")
	assert.Equal(t, "boom", p.GetReason().Error())

	pNil := NewPromise(func(res, rej func(any, error)) { rej(nil, errors.New("boom")) })
	assert.Nil(t, pNil.GetValue(), "reject(nil, err) 时 GetValue 应为 nil")
}

// TestPromise_CombinatorNilHandlerInputs 覆盖组合器接收 NewPromise(nil) 输入的场景：
// 该类 Promise 出生即 Rejected、不经 settle()，其 settled 探测标志在构造期同步置位
// （Batch P1' 免锁快路径的易漏点）。本测试钉住构造期 Store 不被遗漏，
// 且 subscribeDirect/subscribeDirectIndexed 的免锁快路径对四种组合器均正确分发。
func TestPromise_CombinatorNilHandlerInputs(t *testing.T) {
	const nilHandlerReason = "promise handler cannot be nil"

	// 直接钉住构造期置位：缺失时行为断言会被锁路径兜底掩盖（性能静默退化）。
	assert.True(t, NewPromise(nil).settled.Load(), "出生即 Rejected 的 Promise 必须构造期置位 settled")

	t.Run("All should reject with nil-handler reason", func(t *testing.T) {
		ok := NewPromise(func(res, rej func(any, error)) { res(1, nil) })
		p := All(ok, NewPromise(nil))
		assert.Equal(t, Rejected, p.getState())
		assert.Equal(t, nilHandlerReason, p.GetReason().Error())
	})

	t.Run("Race should reject with nil-handler reason", func(t *testing.T) {
		pending := NewPromise(func(res, rej func(any, error)) {})
		p := Race(pending, NewPromise(nil))
		assert.Equal(t, Rejected, p.getState())
		assert.Equal(t, nilHandlerReason, p.GetReason().Error())
	})

	t.Run("AllSettled should record nil-handler error at index", func(t *testing.T) {
		ok := NewPromise(func(res, rej func(any, error)) { res("v", nil) })
		p := AllSettled(NewPromise(nil), ok)
		assert.Equal(t, Fulfilled, p.getState())
		values := p.GetValue().([]any)
		assert.Equal(t, nilHandlerReason, values[0].(error).Error())
		assert.Equal(t, "v", values[1])
	})

	t.Run("Any all nil-handler inputs aggregate reasons", func(t *testing.T) {
		p := Any(NewPromise(nil), NewPromise(nil))
		assert.Equal(t, Rejected, p.getState())
		var agg *AggregateError
		assert.True(t, errors.As(p.GetReason(), &agg), "拒绝原因应为 *AggregateError")
		assert.Len(t, agg.Errors, 2)
		assert.Equal(t, nilHandlerReason, agg.Errors[0].Error())
		assert.Equal(t, nilHandlerReason, agg.Errors[1].Error())
	})
}
