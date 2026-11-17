package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/JacobBananalDev/tavrix/internal/cli"
)

var version = "0.1.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, os.Args[1:], version, os.Stdin, os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}
