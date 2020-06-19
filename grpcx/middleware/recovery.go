package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/msales/pkg/v4/log"
	"google.golang.org/grpc"
)

// RecoveryFunc is used to configure the recovery interceptors.
type RecoveryFunc func(*recoveryConfig)

// WithoutStack disables the stack trace dump from the recovery log.
func WithoutStack() RecoveryFunc {
	return func(cfg *recoveryConfig) {
		cfg.withStack = false
	}
}

// recoveryConfig represents the configuration of the recovery interceptors.
type recoveryConfig struct {
	withStack bool
}

// newRecoveryConfig returns a new config object with sane defaults.
func newRecoveryConfig(opts ...RecoveryFunc) *recoveryConfig {
	cfg := &recoveryConfig{withStack: true}
	cfg.applyOpts(opts)

	return cfg
}

func (cfg *recoveryConfig) applyOpts(opts []RecoveryFunc) {
	for _, fn := range opts {
		fn(cfg)
	}
}

// WithUnaryServerRecovery returns an interceptor that recovers from panics.
func WithUnaryClientRecovery(recoveryOpts ...RecoveryFunc) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		cfg := newRecoveryConfig(recoveryOpts...)

		defer recoveryFunc(ctx, cfg.withStack)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// WithUnaryServerRecovery returns an interceptor that recovers from panics.
func WithUnaryServerRecovery(opts ...RecoveryFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		cfg := newRecoveryConfig(opts...)

		defer recoveryFunc(ctx, cfg.withStack)

		return handler(ctx, req)
	}
}

// WithStreamServerRecovery returns an interceptor that recovers from panics.
func WithStreamServerRecovery(opts ...RecoveryFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cfg := newRecoveryConfig(opts...)

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
