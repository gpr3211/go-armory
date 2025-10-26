package monad

import (
	"sync"
	"time"
)

type Future[A any] struct {
	value    A             // computed value of the Future.
	err      error         // The error from the computation, if any.
	done     bool          // Indicates if the computation is completed.
	mutex    sync.RWMutex  // Mutex to protect access to `value` and `err`.
	waitChan chan struct{} // Channel to signal when the computation is complete.
}

// Creates a new Future from a computation
func NewFuture[A any](compute func() (A, error)) *Future[A] {
	f := &Future[A]{
		waitChan: make(chan struct{}), // close when done
	}
	go func() {
		value, err := compute()
		f.mutex.Lock()
		f.value = value
		f.err = err
		f.done = true
		f.mutex.Unlock()
		close(f.waitChan) // Signal done
	}()

	return f
}

// Example usage of NewFuture:
// future := NewFuture(func() (int, error) {
//     time.Sleep(1 * time.Second)
//     return 42, nil
// })

// Map applies a transformation function to the Future's result, returning a new Future.
func Map[A, B any](f *Future[A], fn func(A) B) *Future[B] {
	return NewFuture(func() (B, error) {
		value, err := f.Get() // Wait for `f` to complete
		if err != nil {
			return *new(B), err
		}
		return fn(value), nil // Apply transformation and return result
	})
}

// Example usage of Map:
// doubled := Map(future, func(x int) int { return x * 2 })

// FlatMap chains two Futures, allowing you to use the result of one to start another Future.
func FlatMap[A, B any](f *Future[A], fn func(A) *Future[B]) *Future[B] {
	return NewFuture(func() (B, error) {
		value, err := f.Get() // Wait for `f` to complete
		if err != nil {
			return *new(B), err
		}
		return fn(value).Get() // Execute the next Future and return its result
	})
}

// Example usage of FlatMap:
// result := FlatMap(future, func(x int) *Future[string] {
//     return NewFuture(func() (string, error) {
//         return fmt.Sprintf("Result is %d", x), nil
//     })
// })

// Get waits for the Future to complete and returns its result or error.
func (f *Future[A]) Get() (A, error) {
	<-f.waitChan    // Block until `waitChan` is closed
	f.mutex.RLock() // Lock to read `value` and `err`
	defer f.mutex.RUnlock()
	return f.value, f.err // Return computed value or error
}

// Example usage of Get:
// result, err := future.Get()

// GetWithTimeout waits for the Future to complete or times out after `timeout`.
// computational error returned by error return, bool returns true if timeout is exceeded.
func (f *Future[A]) GetWithTimeout(timeout time.Duration) (A, error, bool) {
	select {
	case <-f.waitChan: // Future completed
		f.mutex.RLock()
		defer f.mutex.RUnlock()
		return f.value, f.err, false
	case <-time.After(timeout): // Timeout reached
		return *new(A), nil, true
	}
}

// Example usage of GetWithTimeout:
// result, err := future.GetWithTimeout(2 * time.Second)

// Successful returns a Future that completes immediately with a successful value.
func Successful[A any](value A) *Future[A] {
	return NewFuture(func() (A, error) {
		return value, nil
	})
}

// Example usage of Successful:
// successFuture := Successful(100)

// Failed returns a Future that completes immediately with an error.
func Failed[A any](err error) *Future[A] {
	return NewFuture(func() (A, error) {
		return *new(A), err
	})
}

// Example usage of Failed:
// errorFuture := Failed[int](fmt.Errorf("an error occurred"))

// Sequence takes a slice of Futures and returns a Future of a slice of their results.
func Sequence[A any](futures ...*Future[A]) *Future[[]A] {
	return NewFuture(func() ([]A, error) {
		results := make([]A, len(futures))
		for i, future := range futures {
			value, err := future.Get()
			if err != nil {
				return nil, err
			}
			results[i] = value
		}
		return results, nil
	})
}

