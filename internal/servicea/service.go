package servicea

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

// Item item data
type Item struct {
	ID        string
	Code      string
	Name      string
	Signature string
}

// Service service a
type Service struct {
	serviceb *ServiceBClient
}

// NewService create new service
func NewService(serviceb *ServiceBClient) *Service {
	service := Service{
		serviceb: serviceb,
	}
	return &service
}

// GetItem get the item with the given id.
func (s *Service) GetItem(ctx context.Context, id string) (*Item, error) {
	log.Println("msg", "servicea.GetItem")
	if id == "" {
		return nil, errors.New("must provide a valid item id")
	}
	item, err := s.serviceb.GetItem(ctx, id)
	if err != nil {
		log.Println("level", "ERROR", "msg", "item cannot be found because there was an error", "error", err)
		return nil, errors.New("item cannot be found because there was an error")
	}
	if item == nil {
		return nil, nil
	}
	itemFound := Item{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Signature: uuid.New().String(),
	}
	return &itemFound, nil
}
