package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/msales/pkg/v4/grpcx/middleware"
	"github.com/msales/pkg/v4/log"
	"github.com/msales/pkg/v4/stats"
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
