package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yurkovlyanets/gogo/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := server.New()
	if err != nil {
		log.Fatalf("init error: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("run error: %v", err)
	}
}
