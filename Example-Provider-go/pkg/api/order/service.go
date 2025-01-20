package order

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type logger interface {
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

type Config struct{}

type Service struct {
	config *Config
	logger logger
}

func NewService(cfg *Config, logger logger) (*Service, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	return &Service{
		config: cfg,
		logger: logger,
	}, nil
}

// /
// GET /<key>
func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	exists := true

	if exists {
		b := fmt.Sprintf("%v", "value")
		w.Write([]byte(b))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// /
// POST /
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, r.Body); err != nil {
		s.logger.Error("fail reading body on create order: %w", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Created order with id %s", "id")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}
