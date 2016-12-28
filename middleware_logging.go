package kate

import (
	"time"

	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
)

func Logging(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		var start = time.Now()

		log.Info(ctx, "started", "remote", r.RemoteAddr, "method", r.Method, "url", r.RequestURI, "body", string(r.RawBody))

		h.ServeHTTP(ctx, w, r)

		log.Info(ctx, "completed", "status_code", w.StatusCode(), "body", string(w.RawBody()), "duration_ms", int64(time.Since(start)/time.Millisecond))
	}
	return ContextHandlerFunc(f)
}
