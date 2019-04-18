package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/msales/pkg/v3/log"
	"google.golang.org/grpc"
)

// recoveryConfig represents the configuration of the recovery interceptors.
type recoveryConfig struct {
	withStack bool
}

// WithUnaryServerRecovery returns an interceptor that recovers from panics.
func WithUnaryServerRecovery(opts ...func(*recoveryConfig)) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		cfg := &recoveryConfig{withStack: true}
		for _, fn := range opts {
			fn(cfg)
		}

		defer recoveryFunc(ctx, cfg.withStack)

		return handler(ctx, req)
	}
}

// WithStreamServerRecovery returns an interceptor that recovers from panics.
func WithStreamServerRecovery(opts ...func(*recoveryConfig)) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cfg := &recoveryConfig{withStack: true}
		for _, fn := range opts {
			fn(cfg)
		}

		defer recoveryFunc(ss.Context(), cfg.withStack)

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
func WithoutStack() func(*recoveryConfig) {
	return func(cfg *recoveryConfig) {
		cfg.withStack = false
	}
}
