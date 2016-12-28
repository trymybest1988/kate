package main

import (
	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
)

func main() {
	mainCtx := context.Background("main")
	log.Debug(mainCtx, "debug message")
	log.Error(mainCtx, "error message", "error", "not found")

	req1Ctx, _ := context.WithCancelAndLogContext(mainCtx, mainCtx.LogContext().With("session", 1))
	log.Info(req1Ctx, "request started")

	sub1Ctx, _ := context.WithCancel(req1Ctx)
	log.Info(sub1Ctx, "handling request, doing sub task")

	log.Info(req1Ctx, "request completed")

	req2Ctx, _ := context.WithCancelAndLogContext(mainCtx, mainCtx.LogContext().With("session", 2))
	log.Info(req2Ctx, "request started")

	sub2Ctx, _ := context.WithCancel(req2Ctx)
	log.Info(sub2Ctx, "handling request, doing sub task")

	log.Info(req2Ctx, "request completed")
}
