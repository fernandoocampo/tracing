package servicec

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/fernandoocampo/tracing/internal/items"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// ItemResponse item data
type ItemResponse struct {
	Code string `json:"code"`
}

type gRPCServer struct {
	getItemCode grpctransport.Handler
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpoints Endpoints) pb.ItemCodeServiceServer {
	return &gRPCServer{
		getItemCode: grpctransport.NewServer(
			endpoints.GetItemCodeEndpoint,
			decodeGetItemCodeRequest,
			encodeGetItemCodeResponse,
		),
	}
}

func (s *gRPCServer) GetItemCode(ctx context.Context, req *pb.ItemCodeRequest) (*pb.ItemCodeResponse, error) {
	_, resp, err := s.getItemCode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*pb.ItemCodeResponse)
	if !ok {
		log.Println("level", "error", "msg", "item code response not valid", fmt.Sprintf("%v", resp))
		return nil, errors.New("cannot query item code")
	}
	return result, nil
}

func decodeGetItemCodeRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req, ok := r.(*pb.ItemCodeRequest)
	if !ok {
		log.Println("level", "error", "method", "decodeGetItemCodeRequest", "invalid item code request", fmt.Sprintf("%+v", r))
		return nil, errors.New("item ID was not provided")
	}
	return req.ItemID, nil
}

func encodeGetItemCodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	var itemCodeResponse pb.ItemCodeResponse
	if response == "" {
		return &itemCodeResponse, nil
	}
	result, ok := response.(string)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot build get item response", "response", response)
		return &itemCodeResponse, errors.New("cannot build get item code response")
	}
	log.Println("level", "INFO", "msg", "item code response", "response", result)
	itemCodeResponse.Code = result
	return &itemCodeResponse, nil
}
