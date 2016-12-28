package retry

import (
	"time"

	"github.com/k81/kate/context"
)

// Set these to override how now is discovered and how sleeping is done
// This is mostly useful for testing, but you never know
var (
	TimeFunc  = time.Now
	SleepFunc = time.Sleep
)

// This is the main interface around which this library is
// built.  It defines a very simple interface for abstracting retry
// logic in your application.
type Strategy interface {
	Next() bool
	HasNext() bool
}

// Retry strategy expanded with reset functionality.
type ResettableStrategy interface {
	Strategy
	Reset()
}

// Useful helper method.  Calls action until it returns true or
// the retry strategy returns false.
func Do(ctx context.Context, strategy Strategy, action func() bool) bool {
	var (
		newctx, cancel = context.WithCancel(ctx)
		result         = make(chan bool, 1)
	)

	defer cancel()

	for strategy.Next() {
		go func() {
			result <- action()
		}()

		select {
		case <-newctx.Done():
			return false
		case succ := <-result:
			if succ {
				return true
			}
		}
	}

	return false
}

func DoWithReset(ctx context.Context, strategy ResettableStrategy, action func() bool) bool {
	var (
		newctx, cancel = context.WithCancel(ctx)
		result         = make(chan bool, 1)
	)

	defer strategy.Reset()
	defer cancel()

	for strategy.Next() {
		go func() {
			result <- action()
		}()

		select {
		case <-newctx.Done():
			return false
		case succ := <-result:
			if succ {
				return true
			}
		}
	}

	return false
}
