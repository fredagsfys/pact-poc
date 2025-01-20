package service

import (
	"context"
	"fmt"

	"consumer.example.com/m/internal/logger"
	"consumer.example.com/m/internal/server"
	"consumer.example.com/m/pkg/api/order"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Logger   *logger.Config `yaml:"logger"`
	Server   *server.Config `yaml:"server"`
	OrderApi *order.Config  `yaml:"test_api"`
}

type Service struct {
	logger *logger.Logger
	server *server.Server
}

func New(ctx context.Context, config *Config) (*Service, error) {
	logger, err := logger.New(config.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	server, err := server.New(config.Server, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}

	orderService, err := order.NewService(config.OrderApi, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create order api: %w", err)
	}

	server.RegisterService(orderService)

	return &Service{
		logger: logger,
		server: server,
	}, nil
}

// ListenAndServe starts the service.
func (s *Service) ListenAndServe(ctx context.Context) error {
	g, errCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.server.ListenAndServe(errCtx); err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
		}
		return nil
	})
	return g.Wait()
}

// Shutdown gracefully stops the service.
func (s *Service) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http: %w", err)
	}
	return nil
}
