package middleware_test

import (
	"testing"

	"github.com/msales/pkg/v3/grpcx/middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithUnaryClientInterceptors(t *testing.T) {
	interceptor := middleware.WithUnaryClientInterceptors()

	assert.Implements(t, (*grpc.DialOption)(nil), interceptor)
}

func TestWithStreamClientInterceptors(t *testing.T) {
	interceptor := middleware.WithStreamClientInterceptors()

	assert.Implements(t, (*grpc.DialOption)(nil), interceptor)
}

func TestWithUnaryServerInterceptors(t *testing.T) {
	interceptor := middleware.WithUnaryServerInterceptors()
	_, ok := interceptor.(grpc.ServerOption)

	assert.True(t, ok)
}

func TestWithStreamServerInterceptors(t *testing.T) {
	interceptor := middleware.WithStreamServerInterceptors()
	_, ok := interceptor.(grpc.ServerOption)

	assert.True(t, ok)
}
