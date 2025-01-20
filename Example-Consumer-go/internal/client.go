package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Items      []Item    `json:"items"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *OrderClient {
	return &OrderClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *OrderClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	url := fmt.Sprintf("%s/orders/%s", c.baseURL, orderID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &order, nil
}

func (c *OrderClient) CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	url := fmt.Sprintf("%s/orders", c.baseURL)
	body, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("marshaling order: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdOrder Order
	if err := json.NewDecoder(resp.Body).Decode(&createdOrder); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdOrder, nil
}
