package middleware

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"google.golang.org/grpc"
)

// WithUnaryServerLogger adds the logger instance to the unary request context.
func WithUnaryServerLogger(l log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(log.WithLogger(ctx, l), req)
	}
}

// WithStreamServerLogger adds the stats instance to the stream context.
func WithStreamServerLogger(l log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ws := grpc_middleware.WrapServerStream(ss)
		ws.WrappedContext = log.WithLogger(ss.Context(), l)

		return handler(srv, ws)
	}
}

// WithUnaryServerStats adds the stats instance to the unary request context.
func WithUnaryServerStats(s stats.Stats) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(stats.WithStats(ctx, s), req)
	}
}

// WithStreamServerStats adds the stats instance to the stream context.
func WithStreamServerStats(s stats.Stats) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ws := grpc_middleware.WrapServerStream(ss)
		ws.WrappedContext = stats.WithStats(ss.Context(), s)

		return handler(srv, ws)
	}
}

// WithUnaryClientContextTimeout adds timeout to unary client request context.
func WithUnaryClientContextTimeout(d time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// WithUnaryClientLogger adds the logger instance to the unary client request context.
func WithUnaryClientLogger(l log.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(log.WithLogger(ctx, l), method, req, reply, cc, opts...)
	}
}

// WithUnaryClientStats adds the stats instance to the unary client request context.
func WithUnaryClientStats(s stats.Stats) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(stats.WithStats(ctx, s), method, req, reply, cc, opts...)
	}
}
