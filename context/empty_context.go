package context

import (
	"time"

	"github.com/k81/kate/log"
)

type emptyCtx struct {
	name   string
	logCtx *log.LogContext
}

func newEmptyCtx(name string, logCtx *log.LogContext) Context {
	c := &emptyCtx{
		name:   name,
		logCtx: logCtx,
	}
	return c
}

func (c *emptyCtx) String() string {
	return c.name
}

func (c *emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *emptyCtx) Done() <-chan struct{} {
	return nil
}

func (c *emptyCtx) Err() error {
	return nil
}

func (c *emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (c *emptyCtx) LogContext() *log.LogContext {
	return c.logCtx
}
