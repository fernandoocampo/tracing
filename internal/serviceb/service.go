package serviceb

import (
	"context"
	"errors"
	"log"
)

// Item item data
type Item struct {
	ID   string
	Code string
	Name string
}

// Service service a
type Service struct {
	servicecClient *ServiceCClient
	db             map[string]Item
}

// NewService create new service
func NewService(serviceCClient *ServiceCClient) *Service {
	newDB := map[string]Item{
		"1": {
			ID:   "1",
			Name: "item one",
		},
		"2": {
			ID:   "2",
			Name: "item two",
		},
		"3": {
			ID:   "3",
			Name: "item three",
		},
		"4": {
			ID:   "4",
			Name: "item four",
		},
		"5": {
			ID:   "5",
			Name: "item five",
		},
		"6": {
			ID:   "6",
			Name: "item six",
		},
	}
	service := Service{
		servicecClient: serviceCClient,
		db:             newDB,
	}
	return &service
}

// GetItem get the item with the given id.
func (s *Service) GetItem(ctx context.Context, id string) (*Item, error) {
	log.Println("msg", "serviceb.GetItem")
	if id == "" {
		return nil, errors.New("must provide a valid item id")
	}
	item, ok := s.db[id]
	if !ok {
		return nil, nil
	}

	itemCode, err := s.servicecClient.GetItemCode(ctx, id)
	if err != nil {
		log.Println("level", "error", "msg", "cannot get code from item", "error", err)
		item.Code = "TBD"
	}
	item.Code = itemCode
	log.Println("msg", "serviceb.GetItem", "item found", item)

	return &item, nil
}

// PostItem get the item with the given id.
func (s *Service) PostItem(ctx context.Context, item Item) error {
	if item.ID == "" {
		return errors.New("must provide a valid item")
	}
	_, ok := s.db[item.ID]
	if ok {
		return errors.New("the item already exists")
	}
	s.db[item.ID] = item
	return nil
}
