package gateway

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

type Endpoint struct {
	Network, Addr string
}

type Options struct {
	Addr       string
	GRPCServer Endpoint
	OpenAPIDir string
	Mux        []runtime.ServeMuxOption
}
