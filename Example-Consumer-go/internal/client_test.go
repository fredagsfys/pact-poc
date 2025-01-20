package client

import (
	"context"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestOrderAPIClient(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "OrderAPIConsumer",
		Provider: "OrderAPI",
	})
	assert.NoError(t, err)

	// Arrange: Setup our expected interactions
	mockProvider.
		AddInteraction().
		Given("A user with ID 10 exists").
		UponReceiving("A request for Order 10").
		WithRequest("GET", "/order/10").
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json")).
				JSONBody(map[string]interface{}{
					"id": "10",
				})
		})

	// Act: test our API client behaves correctly
	err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		client := NewClient(config.Host)
		ctx := context.Background()

		user, err := client.GetOrder(ctx, "10")

		assert.NoError(t, err)
		assert.Equal(t, 10, user.ID)

		return err
	})

	// Assert
	assert.NoError(t, err)
}
