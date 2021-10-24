package serviceb

import (
	"context"
	"errors"
)

// Item item data
type Item struct {
	ID   string
	Code string
	Name string
}

// Service service a
type Service struct {
	db map[string]Item
}

// NewService create new service
func NewService() *Service {
	newDB := map[string]Item{
		"1": {
			ID:   "1",
			Code: "one",
			Name: "item one",
		},
		"2": {
			ID:   "2",
			Code: "two",
			Name: "item two",
		},
		"3": {
			ID:   "3",
			Code: "three",
			Name: "item three",
		},
		"4": {
			ID:   "4",
			Code: "four",
			Name: "item four",
		},
		"5": {
			ID:   "5",
			Code: "five",
			Name: "item five",
		},
		"6": {
			ID:   "6",
			Code: "six",
			Name: "item six",
		},
	}
	service := Service{
		db: newDB,
	}
	return &service
}

// GetItem get the item with the given id.
func (s *Service) GetItem(ctx context.Context, id string) (*Item, error) {
	if id == "" {
		return nil, errors.New("must provide a valid item id")
	}
	item, ok := s.db[id]
	if !ok {
		return nil, nil
	}
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
