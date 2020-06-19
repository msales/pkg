package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msales/pkg/v3/breaker"
	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	. "github.com/msales/pkg/v3/grpcx/middleware"
)

var breakerErr = errors.New("breaker: circuit breaker is open")

func TestWithBreaker(t *testing.T) {
	s := new(mockStats)
	s.AssertNotCalled(t, "Inc")
	ctx := context.Background()
	ctx = stats.WithStats(ctx, s)

	br := breaker.NewBreaker(
		breaker.RateFuse(1),
		breaker.WithSleep(1*time.Second),
		breaker.WithTestRequests(1),
	)
	interceptor := WithClientBreaker(br, "test")
	err := interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	})

	assert.Nil(t, err)
	s.AssertExpectations(t)
}

func TestWithBreaker_Errored(t *testing.T) {
	s := new(mockStats)
	s.On("Inc", "breaker.error", int64(1), float32(1.0), []interface{}{"state", "open", "name", "test"}).Return(nil).Once()
	ctx := context.Background()
	ctx = stats.WithStats(ctx, s)

	br := breaker.NewBreaker(
		breaker.RateFuse(10),
		breaker.WithSleep(1*time.Second),
		breaker.WithTestRequests(1),
	)

	interceptor := WithClientBreaker(br, "test")
	err := interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return testErr
	})

	assert.Equal(t, testErr, err)

	err = interceptor(ctx, "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return testErr
	})

	assert.Equal(t, breakerErr, err)
	s.AssertExpectations(t)
}

type mockStats struct {
	mock.Mock
}

func (m *mockStats) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *mockStats) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *mockStats) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *mockStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *mockStats) Close() error {
	args := m.Called()
	return args.Error(0)
}
