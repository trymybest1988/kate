package context

import (
	"fmt"
	"sync/atomic"

	"github.com/k81/kate/log"
)

func WithValue(parent Context, key interface{}, val interface{}) Context {
	return WithValueAndLogContext(parent, key, val, parent.LogContext())
}

func WithValueAndLogContext(parent Context, key interface{}, val interface{}, logCtx *log.LogContext) Context {
	c := &valueCtx{
		Context: parent,
		key:     key,
		val:     val,
	}
	c.logCtx.Store(logCtx)
	return c
}

type valueCtx struct {
	Context
	logCtx   atomic.Value
	key, val interface{}
}

func (c *valueCtx) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.val)
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

func (c *valueCtx) LogContext() *log.LogContext {
	return c.logCtx.Load().(*log.LogContext)
}
