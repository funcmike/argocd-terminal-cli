package main

import (
	"context"
	"errors"
	"flag"
	"github.com/funcmike/argocd-terminal-cli/internal/cli"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer cancel()

		if err := cli.Run(os.Args[1:], ctx); err != nil {
			if errors.Is(err, flag.ErrHelp) || errors.Is(err, context.Canceled) {
				return
			}
			panic(err)
		}
	}()

	select {
	case <-signals:
		cancel()
	case <-ctx.Done():
	}

	<-done
}
