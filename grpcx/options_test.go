package grpcx_test

import (
	"context"
	"testing"
	"time"

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

func TestUnaryClientCommonOpts(t *testing.T) {
	opts := grpcx.UnaryClientCommonOpts(context.Background(), 1*time.Second)

	assert.Len(t, opts, 3)
}
