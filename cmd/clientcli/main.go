package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const host = "http://localhost:8890"

// ServiceAItem item data
type ServiceAItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func main() {

}

// GetItem gets item data with given id from service b.
func getItem(ctx context.Context, id string) (*ServiceAItem, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	if id == "" {
		return nil, errors.New("must provide a valid id")
	}
	endpoint := fmt.Sprintf("%s/items/%s", host, id)

	newRequest, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(newRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result ServiceAItem

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
