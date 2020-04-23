package grpcx_test

import (
	"context"
	"testing"

	"github.com/msales/pkg/v4/grpcx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/stats"
)

func TestWithHandlers(t *testing.T) {
	h := grpcx.WithHandlers()

	assert.Implements(t, (*stats.Handler)(nil), h)
}

func TestWithRPCStats(t *testing.T) {
	h := grpcx.WithRPCStats(nil)

	assert.Implements(t, (*stats.Handler)(nil), h)
}

func TestAggregateHandler_TagRPC(t *testing.T) {
	h1, h2 := new(mockHandler), new(mockHandler)
	ctx := context.Background()
	info := &stats.RPCTagInfo{}

	h1.On("TagRPC", ctx, info).Return(ctx)
	h2.On("TagRPC", ctx, info).Return(ctx)

	h := grpcx.WithHandlers(h1, h2)

	retCtx := h.TagRPC(ctx, info)

	h1.AssertExpectations(t)
	h2.AssertExpectations(t)
	assert.Equal(t, ctx, retCtx)
}

func TestAggregateHandler_HandleRPC(t *testing.T) {
	h1, h2 := new(mockHandler), new(mockHandler)
	ctx := context.Background()
	s := &stats.Begin{}

	h1.On("HandleRPC", ctx, s).Return(ctx)
	h2.On("HandleRPC", ctx, s).Return(ctx)

	h := grpcx.WithHandlers(h1, h2)

	h.HandleRPC(ctx, s)

	h1.AssertExpectations(t)
	h2.AssertExpectations(t)
}

func TestAggregateHandler_TagConn(t *testing.T) {
	h1, h2 := new(mockHandler), new(mockHandler)
	ctx := context.Background()
	info := &stats.ConnTagInfo{}

	h1.On("TagConn", ctx, info).Return(ctx)
	h2.On("TagConn", ctx, info).Return(ctx)

	h := grpcx.WithHandlers(h1, h2)

	retCtx := h.TagConn(ctx, info)

	h1.AssertExpectations(t)
	h2.AssertExpectations(t)
	assert.Equal(t, ctx, retCtx)
}

func TestAggregateHandler_HandleConn(t *testing.T) {
	h1, h2 := new(mockHandler), new(mockHandler)
	ctx := context.Background()
	s := &stats.ConnBegin{}

	h1.On("HandleConn", ctx, s).Return(ctx)
	h2.On("HandleConn", ctx, s).Return(ctx)

	h := grpcx.WithHandlers(h1, h2)

	h.HandleConn(ctx, s)

	h1.AssertExpectations(t)
	h2.AssertExpectations(t)
}

type mockHandler struct {
	mock.Mock
}

func (h *mockHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	args := h.Called(ctx, info)

	return args.Get(0).(context.Context)
}

func (h *mockHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	h.Called(ctx, s)
}

func (h *mockHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	args := h.Called(ctx, info)

	return args.Get(0).(context.Context)
}

func (h *mockHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	h.Called(ctx, s)
}
