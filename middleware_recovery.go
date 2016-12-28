package kate

import (
	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
)

func Recovery(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case ErrorInfo:
					Error(ctx, w, v)
				default:
					Error(ctx, w, ErrServerInternal)
					log.Error(ctx, "PANIC", "error", v, "stack", string(log.GetStack(2)))
				}
			}
		}()

		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}
