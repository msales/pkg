package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/msales/pkg/v4/grpcx/middleware"
	"github.com/msales/pkg/v4/log"
	"github.com/msales/pkg/v4/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestWithUnaryServerRecovery(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "test", "stack", mock.AnythingOfType("string"))

	interceptor := middleware.WithUnaryServerRecovery()

	ret, err := interceptor(log.WithLogger(ctx, logger), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test")
	})

	assert.Nil(t, ret)
	assert.NoError(t, err)
	logger.AssertExpectations(t)
}

func TestWithUnaryServerRecovery_WithoutStack(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "test")

	interceptor := middleware.WithUnaryServerRecovery(middleware.WithoutStack())

	ret, err := interceptor(log.WithLogger(ctx, logger), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test")
	})

	assert.Nil(t, ret)
	assert.NoError(t, err)
	logger.AssertExpectations(t)
}

func TestWithUnaryServerRecovery_WithError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	interceptor := middleware.WithUnaryServerRecovery()

	ret, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic(errors.New("test: error"))
	})

	assert.Nil(t, ret)
	assert.NoError(t, err)
}

func TestWithStreamServerRecovery(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "test", "stack", mock.AnythingOfType("string"))

	interceptor := middleware.WithStreamServerRecovery()
	stream := &serverStreamMock{ctx: log.WithLogger(ctx, logger)}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("test")
	})

	assert.NoError(t, err)
	logger.AssertExpectations(t)
}

func TestWithStreamServerRecovery_WithoutStack(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "test")

	interceptor := middleware.WithStreamServerRecovery(middleware.WithoutStack())
	stream := &serverStreamMock{ctx: log.WithLogger(ctx, logger)}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("test")
	})

	assert.NoError(t, err)
	logger.AssertExpectations(t)
}

func TestWithStreamServerRecovery_WithError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	interceptor := middleware.WithStreamServerRecovery()
	stream := &serverStreamMock{ctx: context.Background()}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic(errors.New("test: error"))
	})

	assert.NoError(t, err)
}
