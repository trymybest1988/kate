package context

import (
	"errors"
	"time"

	"github.com/k81/kate/log"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)

	Done() <-chan struct{}

	Err() error

	Value(key interface{}) interface{}

	LogContext() *log.LogContext
}

var (
	Canceled         = errors.New("context canceled")
	DeadlineExceeded = errors.New("context deadline exceeded")
)

func Background(name string) Context {
	return newEmptyCtx(name, log.RootContext.With("module", name))
}

func TODO() Context {
	return newEmptyCtx("TODO", log.RootContext)
}
