package tracers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
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

// HTTPServerTrace enables OpenCensus tracing of a Go kit HTTP transport server.
func HTTPServerTrace() kithttp.ServerOption {
	tracingValues := NewTracingValues()
	serverBefore := kithttp.ServerBefore(
		func(ctx context.Context, req *http.Request) context.Context {
			// Loop over header names
			for name, values := range req.Header {
				// Loop over all values for the name.
				for _, value := range values {
					fmt.Println("header", name, value)
					if _, ok := tracingValues.m[name]; ok {
						tracingValues.m[name] = value
					}
				}
			}
			ctx = context.WithValue(ctx, TracingHeadersKey, tracingValues)
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

// PopulateHeaders populate headers
func PopulateHeaders(ctx context.Context, request *http.Request) *http.Request {
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
