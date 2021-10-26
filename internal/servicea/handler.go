package servicea

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/fernandoocampo/tracing/internal/tracers"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// ItemResponse item data
type ItemResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

// NewHTTPServer is a factory to create http servers for this project.
func NewHTTPServer(endpoints Endpoints) http.Handler {
	router := mux.NewRouter()
	// Add the GO kit HTTP transport middleware to our serverOptions.
	options := []httptransport.ServerOption{
		tracers.HTTPServerTrace(),
		httptransport.ServerErrorEncoder(nil),
	}
	router.Methods(http.MethodGet).Path("/items/{id}").Handler(
		httptransport.NewServer(
			endpoints.GetItemEndpoint,
			decodeGetItemWithIDRequest,
			encodeGetItemWithIDResponse,
			options...,
		),
	)
	return router
}

func decodeGetItemWithIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	userIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("item ID was not provided")
	}
	return userIDParam, nil
}

func encodeGetItemWithIDResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return json.NewEncoder(w).Encode(nil)
	}
	result, ok := response.(*Item)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot build get item response", "response", response)
		return errors.New("cannot build get item response")
	}
	if result == nil {
		return nil
	}
	itemResult := ItemResponse{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		Signature: result.Signature,
	}
	return json.NewEncoder(w).Encode(itemResult)
}
