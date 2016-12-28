package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k81/kate"
	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
)

var (
	mctx = context.Background("http")
)

type HelloWorld struct{}

func (h *HelloWorld) ServeHTTP(ctx context.Context, w kate.ResponseWriter, r *kate.Request) {
	kate.Ok(ctx, w)
}

type NotFound struct{}

func (h *NotFound) ServeHTTP(ctx context.Context, w kate.ResponseWriter, r *kate.Request) {
	panic(kate.ErrNotImplemented)
}

func main() {
	chain := kate.NewChain(
		kate.Logging,
		kate.Timeout(3*time.Second),
		kate.Recovery,
	)

	router := kate.NewRouter(mctx)
	router.GET(
		"/hello",
		chain.Then(&HelloWorld{}),
	)
	router.SetNotFound(chain.Then(&NotFound{}))

	service := kate.Service(mctx, router, nil)
	service.Listen(":8080")

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Info(mctx, "signal got, shutdown ...", "signal", <-sigCh)

	service.Stop()
}
