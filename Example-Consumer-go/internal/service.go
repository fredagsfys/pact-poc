package internal

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	server *server.Server
}

func New(ctx context.Context) (*Service, error) {
	orderService, err := order.NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create monster api: %w", err)
	}

	server.RegisterService(orderService)

	return &Service{
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
