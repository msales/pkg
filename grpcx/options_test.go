package grpcx_test

import (
	"context"
	"testing"

	"github.com/msales/pkg/v3/grpcx"
	"github.com/stretchr/testify/assert"
)

func TestUnaryServerCommonOpts(t *testing.T) {
	opts := grpcx.UnaryServerCommonOpts(context.Background())

	assert.Len(t, opts, 2)
}

func TestStreamServerCommonOpts(t *testing.T) {
	opts := grpcx.StreamServerCommonOpts(context.Background())

	assert.Len(t, opts, 2)
}
