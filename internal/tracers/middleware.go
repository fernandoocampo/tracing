package tracers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type key int

// TracingHeadersKey context tracing headers key
const TracingHeadersKey key = iota

// Tracing keys
const (
	RequestID    = "x-request-id"
	TraceID      = "x-b3-traceid"
	SpanID       = "x-b3-spanid"
	ParentSpanID = "x-b3-parentspanid"
	Samped       = "x-b3-sampled"
	Flags        = "x-b3-flags"
	SpanContext  = "x-ot-span-context"
)

// headers Tracer headers
var headers = [7]string{RequestID, TraceID, SpanID, ParentSpanID, Samped, Flags, SpanContext}

// TracingValues contains values for tracing
type TracingValues struct {
	m map[string]string
}

// NewTracingValues create a new tracing values instance.
func NewTracingValues() *TracingValues {
	newM := make(map[string]string)
	for _, v := range headers {
		newM[v] = ""
	}
	newTracingValues := TracingValues{
		m: newM,
	}
	return &newTracingValues
}

// NewTracingValuesWithData create a new tracing values instance using the given parameters.
func NewTracingValuesWithData(tracingValues map[string]string) *TracingValues {
	var newTracingValues TracingValues
	if len(tracingValues) == 0 {
		return NewTracingValues()
	}
	newTracingValues.m = tracingValues
	return &newTracingValues
}

// HTTPServerTrace enables tracing of a Go kit HTTP transport server.
func HTTPServerTrace() kithttp.ServerOption {
	serverBefore := kithttp.ServerBefore(
		func(ctx context.Context, req *http.Request) context.Context {
			tracingValues := ReadIncomingHTTPHeaders(req)
			if tracingValues == nil {
				return ctx
			}
			ctx = context.WithValue(ctx, TracingHeadersKey, tracingValues)
			return ctx
		},
	)

	return func(s *kithttp.Server) {
		serverBefore(s)
	}
}

// GRPCServerTrace enables tracing of a Go kit GRPC transport server.
func GRPCServerTrace() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Println("level", "INFO", "msg", "tracers.GRPCServerTrace")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("level", "INFO", "msg", "no metadata")
			return handler(ctx, req)
		}

		tracingValuesMap := make(map[string]string)
		for k, v := range md {
			log.Println("level", "INFO", "msg", "metadata", k, v)
			if len(v) == 0 {
				continue
			}
			tracingValuesMap[k] = v[0]
		}
		log.Println("level", "INFO", "tracingValuesMap", tracingValuesMap)

		if len(tracingValuesMap) == 0 {
			log.Println("level", "INFO", "msg", "no metadata")
			return handler(ctx, req)
		}

		tracingValues := NewTracingValuesWithData(tracingValuesMap)
		ctx = context.WithValue(ctx, TracingHeadersKey, tracingValues)

		return handler(ctx, req)
	}
}

// GRPCClientTrace create a grpc UnaryClientInterceptor
func GRPCClientTrace() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		tracingValues := ReadTracingHeadersFromContext(ctx)
		if tracingValues == nil || len(tracingValues.m) == 0 {
			return invoker(ctx, method, req, resp, cc, opts...)
		}
		md := metadata.New(tracingValues.m)
		log.Println("building md", md)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}

// ReadIncomingHTTPHeaders read tracing headers from the given request. Return nil if
// there are not any tracing headers.
func ReadIncomingHTTPHeaders(req *http.Request) *TracingValues {
	tracingValues := NewTracingValues()
	// Loop over header names
	for name, values := range req.Header {
		headerKey := strings.ToLower(name)
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println("header:", headerKey, "value:", value)
			if _, ok := tracingValues.m[headerKey]; ok {
				tracingValues.m[headerKey] = value
			}
		}
	}
	fmt.Println("tracing http headers:", tracingValues.m)
	if len(tracingValues.m) == 0 {
		return nil
	}

	return tracingValues
}

// ReadTracingHeadersFromContext read tracing headers from the given context. Return nil if
// there are not any tracing headers.
func ReadTracingHeadersFromContext(ctx context.Context) *TracingValues {
	ctxValue := ctx.Value(TracingHeadersKey)
	tracingValues, ok := ctxValue.(*TracingValues)

	if !ok {
		log.Println("msg", "invalid tracing headers key", "got", ctxValue)
		return nil
	}
	return tracingValues
}

// PopulateOutgoingHeaders populate outgoing headers, when the service call another service
func PopulateOutgoingHeaders(ctx context.Context, request *http.Request) *http.Request {
	if request == nil {
		return nil
	}
	ctxValue := ctx.Value(TracingHeadersKey)
	tracingHeadersKey, ok := ctxValue.(*TracingValues)

	if !ok {
		log.Println("msg", "invalid tracing headers key", "got", ctxValue)
	}

	for k, v := range tracingHeadersKey.m {
		log.Println("sending header", k, "with value", v)
		request.Header.Set(k, v)
	}

	return request
}
