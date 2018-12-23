
package middleware_test

import (
	"github.com/msales/pkg/v3/grpcx/middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
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

	assert.IsType(t, (grpc.ServerOption)(nil), interceptor)
}

func TestWithStreamServerInterceptors(t *testing.T) {
	interceptor := middleware.WithStreamServerInterceptors()

	assert.IsType(t, (grpc.ServerOption)(nil), interceptor)
}
