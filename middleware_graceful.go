package kate

import (
	"sync"

	"github.com/k81/kate/context"
)

func Graceful(wg *sync.WaitGroup) Middleware {
	return func(h ContextHandler) ContextHandler {
		f := func(ctx context.Context, w ResponseWriter, r *Request) {
			wg.Add(1)
			defer wg.Done()

			h.ServeHTTP(ctx, w, r)
		}
		return ContextHandlerFunc(f)
	}
}
