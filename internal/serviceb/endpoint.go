package serviceb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints servicea endpoints
type Endpoints struct {
	GetItemEndpoint endpoint.Endpoint
}

// NewEndpoints create service a endpoints
func NewEndpoints(service *Service) Endpoints {
	return Endpoints{
		GetItemEndpoint: makeGetItemEndpoint(service),
	}
}

func makeGetItemEndpoint(service *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		itemID, ok := request.(string)
		if !ok {
			log.Println("level", "error", "invalid get item request", fmt.Sprintf("%+v", request))
			return nil, errors.New("invalid request")
		}
		item, err := service.GetItem(ctx, itemID)
		if err != nil {
			return nil, nil
		}
		return item, nil
	}
}
