package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/msales/pkg/v3/grpcx/middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithUnaryServerRecovery(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	interceptor := middleware.WithUnaryServerRecovery()

	ret, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test")
	})

	assert.Nil(t, ret)
	assert.NoError(t, err)
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

	interceptor := middleware.WithStreamServerRecovery()
	stream := &serverStreamMock{ctx: context.Background()}

	err := interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("test")
	})

	assert.NoError(t, err)
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
