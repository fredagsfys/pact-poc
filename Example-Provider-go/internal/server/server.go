package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

type Config struct {
	Addr string `yaml:"addr"`
	TLS  *struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
}

// Server is a http server.
type Server struct {
	server http.Server
	logger logger
	Router *mux.Router
	config *Config
}

type orderService interface {
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

func New(cfg *Config, logger logger) (*Server, error) {
	router := mux.NewRouter()
	router.Use(commonMiddleware)

	return &Server{
		logger: logger,
		Router: router,
		config: cfg,
	}, nil
}

func (s *Server) RegisterService(orderService orderService) {
	s.Router.HandleFunc("/api/order/{id}", orderService.Get).Methods("GET")
	s.Router.HandleFunc("/api/order", orderService.Create).Methods("POST")
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// ListenAndServe starts http server.
func (s *Server) ListenAndServe(ctx context.Context) error {
	addr := ":8080"
	if s.config.Addr != "" {
		addr = s.config.Addr
	}

	s.server = http.Server{
		Addr:         addr,
		Handler:      s.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if s.config.TLS == nil {
		s.logger.Info("serving http on %s", addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	}

	s.logger.Info("serving https on %s", addr)
	if err := s.server.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown gracefully stops http server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down http server")
	s.server.Shutdown(ctx)

	return nil
}
