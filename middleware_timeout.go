package kate

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/k81/kate/context"
)

func Timeout(timeout time.Duration) Middleware {
	return func(h ContextHandler) ContextHandler {
		f := func(ctx context.Context, w ResponseWriter, r *Request) {
			var (
				cancel context.CancelFunc
				done   = make(chan struct{})
			)

			if timeout <= 0 {
				h.ServeHTTP(ctx, w, r)
			} else {
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()

				tw := &timeoutResponseWriter{
					ResponseWriter: w,
					h:              make(http.Header),
				}

				go func() {
					h.ServeHTTP(ctx, tw, r)
					close(done)
				}()

				select {
				case <-done:
					tw.mu.Lock()
					defer tw.mu.Unlock()
					dst := w.Header()
					for k, vv := range tw.h {
						dst[k] = vv
					}
					w.WriteHeader(tw.code)
					w.Write(tw.wbuf.Bytes())
				case <-ctx.Done():
					tw.mu.Lock()
					defer tw.mu.Unlock()
					Error(ctx, w, ErrTimeout)
					tw.timedOut = true
				}
			}
		}
		return ContextHandlerFunc(f)
	}
}

type timeoutResponseWriter struct {
	ResponseWriter
	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

func (tw *timeoutResponseWriter) Header() http.Header {
	return tw.h
}

func (tw *timeoutResponseWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, ErrTimeout
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func (tw *timeoutResponseWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutResponseWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *timeoutResponseWriter) StatusCode() int {
	return tw.code
}

func (tw *timeoutResponseWriter) RawBody() []byte {
	return tw.wbuf.Bytes()
}

func (tw *timeoutResponseWriter) WriteJson(v interface{}) error {
	b, err := tw.EncodeJson(v)
	if err != nil {
		return err
	}

	tw.SetData(v)
	_, err = tw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (tw *timeoutResponseWriter) Flush() {
	if !tw.wroteHeader {
		tw.WriteHeader(http.StatusOK)
	}
	flusher := tw.ResponseWriter.(http.Flusher)
	flusher.Flush()
}
