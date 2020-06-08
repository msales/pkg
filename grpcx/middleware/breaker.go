package middleware

import (
	"context"
	"errors"

	"github.com/msales/pkg/v3/breaker"
	"github.com/msales/pkg/v3/stats"
	"google.golang.org/grpc"
)

const (
	breakerErrorKey = "breaker.error"
	stateTag        = "state"
)

// WithBreaker adds breaker to client request.
func WithClientBreaker(br *breaker.Breaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := br.Run(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})

		if errors.Is(err, breaker.ErrOpenState) {
			_ = stats.Inc(ctx, breakerErrorKey, 1, 1.0, stateTag, "open")
		}

		if errors.Is(err, breaker.ErrTooManyRequests) {
			_ = stats.Inc(ctx, breakerErrorKey, 1, 1.0, stateTag, "too_many_requests")
		}

		return err
	}
}
