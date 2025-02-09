package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()
				return
			}
		}
	}()

	return ctx
}
