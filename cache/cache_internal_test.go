package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCache(t *testing.T) {
	ctx := WithCache(context.Background(), Null)

	got := ctx.Value(ctxKey)

	assert.Equal(t, Null, got)
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey, Null)

	got, ok := FromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, Null, got)
}

func TestFromContext_NotSet(t *testing.T) {
	ctx := context.Background()

	got, ok := FromContext(ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetCache(t *testing.T) {
	tests := []struct{
		ctx context.Context
		expect Cache
	} {
		{context.Background(), Null},
	}

	for _, tt := range tests {
		got := getCache(tt.ctx)

		assert.Equal(t, tt.expect, got)
	}
}
