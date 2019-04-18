package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/msales/pkg/v3/log"
	"google.golang.org/grpc"
)

// WithUnaryServerRecovery returns an interceptor that recovers from panics.
func WithUnaryServerRecovery(stackOpts ...func() bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		withStack := true
		for _, fn := range stackOpts {
			withStack = fn()
		}

		defer recoveryFunc(ctx, withStack)

		return handler(ctx, req)
	}
}

// WithStreamServerRecovery returns an interceptor that recovers from panics.
func WithStreamServerRecovery(stackOpts ...func() bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		withStack := true
		for _, fn := range stackOpts {
			withStack = fn()
		}

		defer recoveryFunc(ss.Context(), withStack)

		return handler(srv, ss)
	}
}

func recoveryFunc(ctx context.Context, withStack bool) {
	if v := recover(); v != nil {
		err := fmt.Errorf("%v", v)
		if v, ok := v.(error); ok {
			err = v
		}

		var logCtx []interface{}
		if withStack {
			logCtx = append(logCtx, "stack", string(debug.Stack()))
		}

		log.Error(ctx, err.Error(), logCtx...)
	}
}

// WithoutStack disables the stack trace dump from the recovery log.
func WithoutStack() func() bool {
	return func() bool {
		return false
	}
}
