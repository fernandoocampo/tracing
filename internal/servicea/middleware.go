package servicea

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

func loggingMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			log.Println("msg", "calling endpoint")
			defer log.Println("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

// HTTPServerTrace enables OpenCensus tracing of a Go kit HTTP transport server.
func HTTPServerTrace() kithttp.ServerOption {
	serverBefore := kithttp.ServerBefore(
		func(ctx context.Context, req *http.Request) context.Context {
			// Loop over header names
			for name, values := range req.Header {
				// Loop over all values for the name.
				for _, value := range values {
					fmt.Println("header", name, value)
				}
			}
			return ctx
		},
	)

	serverFinalizer := kithttp.ServerFinalizer(
		func(ctx context.Context, code int, r *http.Request) {
		},
	)

	return func(s *kithttp.Server) {
		serverBefore(s)
		serverFinalizer(s)
	}
}
