package order

import (
	"testing"

	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
)

type testLogger struct{}

func (t *testLogger) Warn(format string, v ...interface{}) {}

func (t *testLogger) Error(format string, v ...interface{}) {}

func TestV3HTTPProvider(t *testing.T) {
	// 1. Start your Provider API in the background
	// ctx, cancel := context.WithCancel(context.Background())
	// cfg := &service.Config{}
	// go service.New(ctx, cfg)
	cfg := &Config{}
	logger := &testLogger{}
	go NewService(cfg, logger)

	verifier := provider.NewVerifier()

	// Verify the Provider with pact broker
	// The console will display if the verification was successful or not, the
	// assertions being made and any discrepancies with the contract
	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: "http://localhost:8080",
		BrokerURL:       "localhost:9292/",
	})

	// Ensure the verification succeeded
	assert.NoError(t, err)
}
