package serviceb

import (
	"context"
	"errors"
	"log"

	pb "github.com/fernandoocampo/tracing/internal/items"
	"github.com/fernandoocampo/tracing/internal/tracers"
	"google.golang.org/grpc"
)

// ServiceCClient service c client.
type ServiceCClient struct {
	client pb.ItemCodeServiceClient
}

// NewServicCGRPCClient create a new grpc service c client
func NewServicCGRPCClient() (*ServiceCClient, error) {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(tracers.GRPCClientTrace()))
	if err != nil {
		log.Println("msg", "error dialing StockEventService", "error", err)
		return &ServiceCClient{}, err
	}
	newServiceCClient := ServiceCClient{
		client: pb.NewItemCodeServiceClient(conn),
	}
	return &newServiceCClient, nil
}

// GetItemCode get item code from given item id
func (s *ServiceCClient) GetItemCode(ctx context.Context, itemID string) (string, error) {
	request := pb.ItemCodeRequest{
		ItemID: itemID,
	}

	itemCode, err := s.client.GetItemCode(ctx, &request)
	if err != nil {
		log.Println("msg", "cannot get item code", "error", err)
		return "", errors.New("unexpected error")
	}

	return itemCode.Code, nil
}
