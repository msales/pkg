package grpcx

import (
	"context"
	"time"

	"github.com/msales/pkg/v4/grpcx/middleware"
	"github.com/msales/pkg/v4/log"
	"github.com/msales/pkg/v4/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

// UnaryClientCommonOpts returns commonly options for an unary client.
func UnaryClientCommonOpts(ctx context.Context, timeout time.Duration, additional ...grpc.UnaryClientInterceptor) []grpc.DialOption {
	l, s := getLoggerAndStats(ctx)

	interceptors := []grpc.UnaryClientInterceptor{
		middleware.WithUnaryClientLogger(l),
		middleware.WithUnaryClientStats(s),
		middleware.WithUnaryClientRecovery(),
		middleware.WithUnaryClientContextTimeout(timeout),
	}

	interceptors = append(interceptors, additional...)

	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		middleware.WithUnaryClientInterceptors(interceptors...),
	}
}

// UnaryServerCommonOpts returns commonly options for an unary server.
func UnaryServerCommonOpts(ctx context.Context, statsTagsFns ...TagsFunc) []grpc.ServerOption {
	l, s := getLoggerAndStats(ctx)

	return []grpc.ServerOption{
		middleware.WithUnaryServerInterceptors(
			middleware.WithUnaryServerLogger(l),
			middleware.WithUnaryServerStats(s),
			middleware.WithUnaryServerRecovery(),
		),
		grpc.StatsHandler(
			WithRPCStats(s, statsTagsFns...),
		),
	}
}

// StreamServerCommonOpts returns commonly options for a stream server.
func StreamServerCommonOpts(ctx context.Context, statsTagsFns ...TagsFunc) []grpc.ServerOption {
	l, s := getLoggerAndStats(ctx)

	return []grpc.ServerOption{
		middleware.WithStreamServerInterceptors(
			middleware.WithStreamServerLogger(l),
			middleware.WithStreamServerStats(s),
			middleware.WithStreamServerRecovery(),
		),
		grpc.StatsHandler(
			WithRPCStats(s, statsTagsFns...),
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
