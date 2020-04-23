package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/msales/pkg/v4/log"
	"github.com/msales/pkg/v4/stats"
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
