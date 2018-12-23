package middleware

import (
	"context"
	"fmt"

	"github.com/msales/pkg/v3/log"
	"google.golang.org/grpc"
)

// WithUnaryServerRecovery returns an interceptor that recovers from panics.
func WithUnaryServerRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer recoveryFunc(ctx)

		return handler(ctx, req)
	}
}

// WithStreamServerRecovery returns an interceptor that recovers from panics.
func WithStreamServerRecovery() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer recoveryFunc(ss.Context())

		return handler(srv, ss)
	}
}

func recoveryFunc(ctx context.Context) {
	if v := recover(); v != nil {
		err := fmt.Errorf("%v", v)
		if v, ok := v.(error); ok {
			err = v
		}

		log.Error(ctx, err.Error())
	}
}
