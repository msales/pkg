package grpcx

import (
	"context"

	"github.com/msales/pkg/v3/grpcx/middleware"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"google.golang.org/grpc"
)

// UnaryServerCommonOpts returns commonly options for an unary server.
func UnaryServerCommonOpts(ctx context.Context) []grpc.ServerOption {
	l, s := getLoggerAndStats(ctx)

	return []grpc.ServerOption{
		middleware.WithUnaryServerInterceptors(
			middleware.WithUnaryServerRecovery(),
			middleware.WithUnaryServerLogger(l),
			middleware.WithUnaryServerStats(s),
		),
		grpc.StatsHandler(
			WithRPCStats(s),
		),
	}
}

// StreamServerCommonOpts returns commonly options for a stream server.
func StreamServerCommonOpts(ctx context.Context) []grpc.ServerOption {
	l, s := getLoggerAndStats(ctx)

	return []grpc.ServerOption{
		middleware.WithStreamServerInterceptors(
			middleware.WithStreamServerRecovery(),
			middleware.WithStreamServerLogger(l),
			middleware.WithStreamServerStats(s),
		),
		grpc.StatsHandler(
			WithRPCStats(s),
		),
	}
}

func getLoggerAndStats(ctx context.Context) (log.Logger, stats.Stats) {
	l, ok := log.FromContext(ctx)
	if !ok {
		l = log.Null
	}

	s, ok := stats.FromContext(ctx)
	if !ok {
		s = stats.Null
	}

	return l, s
}
