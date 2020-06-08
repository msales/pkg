package middleware

import (
	"context"

	"github.com/msales/pkg/v4/breaker"
	"github.com/msales/pkg/v4/stats"
	"google.golang.org/grpc"
)

const (
	breakerErrorKey = "breaker.error"
	stateTag        = "state"
)

// WithBreaker adds breaker to client request.
func WithClientBreaker(br *breaker.Breaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		st, ok := stats.FromContext(ctx)
		if !ok {
			st = stats.Null
		}

		err := br.Run(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})

		if err == breaker.ErrOpenState {
			_ = st.Inc(breakerErrorKey, 1, 1.0, stateTag, "open")
		}

		if err == breaker.ErrTooManyRequests {
			_ = st.Inc(breakerErrorKey, 1, 1.0, stateTag, "too_many_requests")
		}

		return err
	}
}
