package servicec

import (
	"context"
	"errors"
	"log"

	"github.com/fernandoocampo/tracing/internal/tracers"
)

// Service service a
type Service struct {
	db map[string]string
}

// NewService create new service
func NewService() *Service {
	newDB := map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
		"4": "four",
		"5": "five",
		"6": "six",
	}
	service := Service{
		db: newDB,
	}
	return &service
}

// GetItemCode get the item code with the given item id.
func (s *Service) GetItemCode(ctx context.Context, itemID string) (string, error) {
	log.Println("msg", "servicec.GetItem")
	tracingValues := tracers.ReadTracingHeadersFromContext(ctx)
	log.Println("msg", "servicec.GetItem", "tracing values", tracingValues)

	if itemID == "" {
		return "", errors.New("must provide a valid item id")
	}
	code, ok := s.db[itemID]
	if !ok {
		return "", nil
	}
	return code, nil
}
