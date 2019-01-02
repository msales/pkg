package grpcx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/stats"
)

func TestRPCStatsHandler_TagRPC(t *testing.T) {
	h := WithRPCStats(nil)
	info := &stats.RPCTagInfo{FullMethodName: "TestMethod"}

	ctx := h.TagRPC(context.Background(), info)

	assert.Equal(t, "TestMethod", ctx.Value(methodKey))
}

func TestRPCStatsHandler_HandleRPC_Begin(t *testing.T) {
	ctx := context.WithValue(context.Background(), methodKey, "TestMethod")
	s := new(mockStats)
	s.On("Inc", "rpc.begin", int64(1), float32(1.0), []interface{}{"method", "TestMethod"}).Return(nil).Once()

	h := WithRPCStats(s)

	h.HandleRPC(ctx, &stats.Begin{})

	s.AssertExpectations(t)
}

func TestRPCStatsHandler_HandleRPC_Begin_WithUnknownMethod(t *testing.T) {
	s := new(mockStats)
	s.On("Inc", "rpc.begin", int64(1), float32(1.0), []interface{}{"method", "unknown"}).Return(nil).Once()

	h := WithRPCStats(s)

	h.HandleRPC(context.Background(), &stats.Begin{})

	s.AssertExpectations(t)
}

func TestRPCStatsHandler_HandleRPC_End(t *testing.T) {
	ctx := context.WithValue(context.Background(), methodKey, "TestMethod")
	s := new(mockStats)
	s.On("Inc", "rpc.end", int64(1), float32(1.0), []interface{}{"method", "TestMethod", "status", "success"}).Return(nil).Once()
	s.On("Timing", "rpc.time", time.Second, float32(1.0), []interface{}{"method", "TestMethod", "status", "success"}).Return(nil).Once()
	now := time.Now()

	h := WithRPCStats(s)

	h.HandleRPC(ctx, &stats.End{BeginTime: now, EndTime: now.Add(time.Second)})

	s.AssertExpectations(t)
}

func TestRPCStatsHandler_HandleRPC_End_WithError(t *testing.T) {
	ctx := context.WithValue(context.Background(), methodKey, "TestMethod")
	s := new(mockStats)
	s.On("Inc", "rpc.end", int64(1), float32(1.0), []interface{}{"method", "TestMethod", "status", "error"}).Return(nil).Once()
	s.On("Timing", "rpc.time", time.Second, float32(1.0), []interface{}{"method", "TestMethod", "status", "error"}).Return(nil).Once()
	now := time.Now()

	h := WithRPCStats(s)

	h.HandleRPC(ctx, &stats.End{BeginTime: now, EndTime: now.Add(time.Second), Error: errors.New("test: error")})

	s.AssertExpectations(t)
}

func TestRPCStatsHandler_HandleRPC_OtherEvent(t *testing.T) {
	ctx := context.WithValue(context.Background(), methodKey, "TestMethod")
	s := new(mockStats)

	h := WithRPCStats(s)

	h.HandleRPC(ctx, &stats.InPayload{})
}

func TestHandler_TagRPC(t *testing.T) {
	h := &handler{}
	ctx := context.Background()
	info := &stats.RPCTagInfo{}

	retCtx := h.TagRPC(ctx, info)

	assert.Equal(t, ctx, retCtx)
	assert.Equal(t, &stats.RPCTagInfo{}, info)
}

func TestHandler_TagConn(t *testing.T) {
	h := &handler{}
	ctx := context.Background()
	info := &stats.ConnTagInfo{}

	retCtx := h.TagConn(ctx, info)

	assert.Equal(t, ctx, retCtx)
	assert.Equal(t, &stats.ConnTagInfo{}, info)
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
