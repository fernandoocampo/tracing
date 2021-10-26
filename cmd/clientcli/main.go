package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	item, err := getItem(ctx, "1")
	if err != nil {
		log.Println("error", err)
		return
	}
	log.Println("item", item)
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
	newRequest.Header.Set("x-request-id", "value-1")
	newRequest.Header.Set("x-b3-traceid", "value-2")
	newRequest.Header.Set("x-b3-spanid", "value-3")
	newRequest.Header.Set("x-b3-parentspanid", "value-4")
	newRequest.Header.Set("x-b3-sampled", "value-5")
	newRequest.Header.Set("x-b3-flags", "value-6")
	newRequest.Header.Set("x-ot-span-context", "value-7")

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
