package kate

import (
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
	"github.com/k81/kate/context"
)

type ContextHandler interface {
	ServeHTTP(context.Context, ResponseWriter, *Request)
}

type ContextHandlerFunc func(context.Context, ResponseWriter, *Request)

func (h ContextHandlerFunc) ServeHTTP(ctx context.Context, w ResponseWriter, r *Request) {
	h(ctx, w, r)
}

var (
	gReqId = uint64(0)
)

func nextReqId() uint64 {
	return atomic.AddUint64(&gReqId, uint64(1))
}

func StdHandler(ctx context.Context, h ContextHandler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var (
			request   *Request
			response  *responseWriter
			newctx    context.Context
			cancel    context.CancelFunc
			requestId uint64 = nextReqId()
			err       error
		)

		newctx, cancel = context.WithCancelAndLogContext(ctx, ctx.LogContext().With("session", requestId))
		defer cancel()

		request = &Request{
			Request: r,
			Env:     map[string]interface{}{},
		}

		response = &responseWriter{
			ResponseWriter: w,
			wroteHeader:    false,
		}

		if err = request.ParseForm(); err != nil {
			Error(ctx, response, ErrBadRequest)
			return
		}

		request.RawBody, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			Error(ctx, response, ErrServerInternal)
			return
		}

		h.ServeHTTP(newctx, response, request)
	}
	return http.HandlerFunc(f)
}

func Handle(ctx context.Context, h ContextHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var (
			request   *Request
			response  *responseWriter
			newctx    context.Context
			cancel    context.CancelFunc
			requestId uint64 = nextReqId()
			err       error
		)

		newctx, cancel = context.WithCancelAndLogContext(ctx, ctx.LogContext().With("session", requestId))
		defer cancel()

		request = &Request{
			Request: r,
			Id:      requestId,
			Path:    params,
			Env:     map[string]interface{}{},
		}

		response = &responseWriter{
			ResponseWriter: w,
			wroteHeader:    false,
		}

		if err = request.ParseForm(); err != nil {
			Error(ctx, response, ErrBadRequest)
			return
		}

		request.RawBody, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			Error(ctx, response, ErrServerInternal)
			return
		}

		h.ServeHTTP(newctx, response, request)
	}
}
