package context

import (
	"fmt"
	"time"

	"github.com/k81/kate/log"
)

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	return WithDeadlineAndLogContext(parent, deadline, parent.LogContext())
}

func WithDeadlineAndLogContext(parent Context, deadline time.Time, logCtx *log.LogContext) (Context, CancelFunc) {
	if cur, ok := parent.Deadline(); ok && cur.Before(deadline) {
		// The current deadline is already sooner than the new one.
		return WithCancelAndLogContext(parent, logCtx)
	}
	c := &timerCtx{
		cancelCtx: cancelCtx{
			Context: parent,
			logCtx:  logCtx,
			done:    make(chan struct{}),
		},
		deadline: deadline,
	}
	propagateCancel(parent, c)
	d := deadline.Sub(time.Now())
	if d <= 0 {
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() { c.cancel(true, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(d, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithTimeoutAndLogContext(parent, timeout, parent.LogContext())
}

func WithTimeoutAndLogContext(parent Context, timeout time.Duration, logCtx *log.LogContext) (Context, CancelFunc) {
	return WithDeadlineAndLogContext(parent, time.Now().Add(timeout), logCtx)
}

type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

func (c *timerCtx) String() string {
	return fmt.Sprintf("%v.WithDeadline(%s [%s])", c.cancelCtx.Context, c.deadline, c.deadline.Sub(time.Now()))
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}
