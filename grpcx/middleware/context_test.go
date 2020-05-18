package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msales/pkg/v3/grpcx/middleware"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var testErr = errors.New("test: error")

func TestWithUnaryServerLogger(t *testing.T) {
	interceptor := middleware.WithUnaryServerLogger(log.Null)

	res, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		l, ok := log.FromContext(ctx)

		assert.Equal(t, l, log.Null)
		assert.True(t, ok)

		return "test", testErr
	})

	assert.Equal(t, "test", res)
	assert.Equal(t, testErr, err)
}

func TestWithStreamServerLogger(t *testing.T) {
	interceptor := middleware.WithStreamServerLogger(log.Null)
	stream := &serverStreamMock{ctx: context.Background()}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		l, ok := log.FromContext(stream.Context())

		assert.Equal(t, l, log.Null)
		assert.True(t, ok)

		return testErr
	})

	assert.Equal(t, testErr, err)
}

func TestWithUnaryServerStats(t *testing.T) {
	interceptor := middleware.WithUnaryServerStats(stats.Null)

	res, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		s, ok := stats.FromContext(ctx)

		assert.Equal(t, s, stats.Null)
		assert.True(t, ok)

		return "test", testErr
	})

	assert.Equal(t, "test", res)
	assert.Equal(t, testErr, err)
}

func TestWithStreamServerStats(t *testing.T) {
	interceptor := middleware.WithStreamServerStats(stats.Null)
	stream := &serverStreamMock{ctx: context.Background()}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		s, ok := stats.FromContext(stream.Context())

		assert.Equal(t, s, stats.Null)
		assert.True(t, ok)

		return testErr
	})

	assert.Equal(t, testErr, err)
}

func TestWithUnaryClientContextTimeout(t *testing.T) {
	ctx := context.Background()

	interceptor := middleware.WithUnaryClientContextTimeout(1 * time.Hour)
	err := interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		_, ok := ctx.Deadline()

		assert.True(t, ok)

		return testErr
	})

	assert.Equal(t, testErr, err)
}

func TestWithUnaryClientLogger(t *testing.T) {
	interceptor := middleware.WithUnaryClientLogger(log.Null)

	err := interceptor(context.Background(), "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		l, ok := log.FromContext(ctx)

		assert.Equal(t, l, log.Null)
		assert.True(t, ok)

		return testErr
	})

	assert.Equal(t, testErr, err)
}

func TestWithUnaryClientStats(t *testing.T) {
	interceptor := middleware.WithUnaryClientStats(stats.Null)

	err := interceptor(context.Background(), "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		s, ok := stats.FromContext(ctx)

		assert.Equal(t, s, stats.Null)
		assert.True(t, ok)

		return testErr
	})

	assert.Equal(t, testErr, err)
}
