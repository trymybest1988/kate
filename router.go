package kate

import (
	"sync"

	"github.com/k81/kate/context"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*sync.WaitGroup
	*httprouter.Router
	ctx context.Context
}

func NewRouter(ctx context.Context) *Router {
	return &Router{
		WaitGroup: &sync.WaitGroup{},
		Router:    httprouter.New(),
		ctx:       ctx,
	}
}

func (r *Router) Handle(method string, path string, h ContextHandler) {
	h = Graceful(r.WaitGroup)(h)
	r.Router.Handle(method, path, Handle(r.ctx, h))
}

func (r *Router) GET(path string, h ContextHandler) {
	r.Handle("GET", path, h)
}

func (r *Router) HEAD(path string, h ContextHandler) {
	r.Handle("HEAD", path, h)
}

func (r *Router) OPTIONS(path string, h ContextHandler) {
	r.Handle("OPTIONS", path, h)
}

func (r *Router) POST(path string, h ContextHandler) {
	r.Handle("POST", path, h)
}

func (r *Router) PUT(path string, h ContextHandler) {
	r.Handle("PUT", path, h)
}

func (r *Router) PATCH(path string, h ContextHandler) {
	r.Handle("PATCH", path, h)
}

func (r *Router) DELETE(path string, h ContextHandler) {
	r.Handle("DELETE", path, h)
}

func (r *Router) SetNotFound(h ContextHandler) {
	r.NotFound = StdHandler(r.ctx, h)
	r.MethodNotAllowed = StdHandler(r.ctx, h)
}
