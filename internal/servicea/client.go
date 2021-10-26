package servicea

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fernandoocampo/tracing/internal/tracers"
)

const (
	serviceBHost = "http://localhost:8087"
)

// ServiceBItem item data
type ServiceBItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ServiceBClient service b client.
type ServiceBClient struct {
	host   string
	client *http.Client
}

// NewServiceBClient service b client
func NewServiceBClient(host string) *ServiceBClient {
	if host == "" {
		host = serviceBHost
	}
	newServiceBClient := ServiceBClient{
		host: host,
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
	return &newServiceBClient
}

// GetItem gets item data with given id from service b.
func (s *ServiceBClient) GetItem(ctx context.Context, id string) (*ServiceBItem, error) {
	if id == "" {
		return nil, errors.New("must provide a valid id")
	}
	endpoint := fmt.Sprintf("%s/items/%s", s.host, id)

	newRequest, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	newRequest = tracers.PopulateOutgoingHeaders(ctx, newRequest)

	response, err := s.client.Do(newRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result ServiceBItem

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
