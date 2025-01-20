package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "consumer.example.com/m/internal"
	"gopkg.in/yaml.v3"
)

func main() {
	path := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *path == "" {
		log.Fatal("No configuration was set, use --config to set config file")
	}

	cfg, err := parseConfig(*path)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	svc, err := service.New(ctx, cfg)
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

func parseConfig(path string) (*service.Config, error) {
	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file '%s': %w", path, err)
	}

	c := &service.Config{}
	if err := yaml.Unmarshal(rawConfig, c); err != nil {
		return nil, fmt.Errorf("failed to parse config as yaml: %w", err)
	}

	return c, nil
}
