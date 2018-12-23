// This file contains simple convenience functions.
// It only chains library function calls, with no added logic.

package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// WithUnaryClientInterceptors wraps multiple unary client interceptors in a single option.
func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
}

// WithStreamClientInterceptors wraps multiple stream client interceptors in a single option.
func WithStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...))
}

// WithUnaryServerInterceptors wraps multiple unary server interceptors in a single option.
func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))
}

// WithStreamServerInterceptors wraps multiple stream server interceptors in a single option.
func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors...))
}
