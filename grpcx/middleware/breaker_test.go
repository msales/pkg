package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msales/pkg/v3/breaker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	. "github.com/msales/pkg/v3/grpcx/middleware"
)

var breakerErr = errors.New("breaker: circuit breaker is open")

func TestWithBreaker(t *testing.T) {
	ctx := context.Background()

	br := breaker.NewBreaker(
		breaker.RateFuse(1),
		breaker.WithSleep(1*time.Second),
		breaker.WithTestRequests(1),
	)

	interceptor := WithClientBreaker(br)
	err := interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	})

	assert.Nil(t, err)
}

func TestWithBreaker_Errored(t *testing.T) {
	ctx := context.Background()

	br := breaker.NewBreaker(
		breaker.RateFuse(10),
		breaker.WithSleep(1*time.Second),
		breaker.WithTestRequests(1),
	)

	interceptor := WithClientBreaker(br)
	err := interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return testErr
	})

	assert.Equal(t, testErr, err)

	err = interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return testErr
	})

	assert.Equal(t, breakerErr, err)
}
