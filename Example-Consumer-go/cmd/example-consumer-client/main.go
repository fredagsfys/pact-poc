package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "consumer.example.com/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	svc, err := service.New(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %s", err)
	}

	// Wait for shut down in a separate goroutine.
	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal, 1)
		signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-shutdownCh

		cancel()

		shutdownTimeout := 30 * time.Second
		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCtxCancel()
		errCh <- svc.Shutdown(shutdownCtx)
	}()

	if err := svc.ListenAndServe(ctx); err != nil {
		log.Fatalf("failed to start service: %s", err)
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		log.Fatal(err.Error())
	}
}
