package stats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStats(t *testing.T) {
	ctx := WithStats(context.Background(), Null)

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

func TestTaggedStats_CollectTags(t *testing.T) {
	stats := &TaggedStats{
		tags: map[string]string{
			"test1": "foo",
			"test2": "bar",
		},
	}

	tests := []struct {
		tags     map[string]string
		expected map[string]string
	}{
		{
			nil,
			map[string]string{"test1": "foo", "test2": "bar"},
		},
		{
			map[string]string{},
			map[string]string{"test1": "foo", "test2": "bar"},
		},
		{
			map[string]string{"foo": "bar", "baz": "bat"},
			map[string]string{"test1": "foo", "test2": "bar", "foo": "bar", "baz": "bat"},
		},
		{
			map[string]string{"foo": "bar", "test1": "bat"},
			map[string]string{"test1": "bat", "test2": "bar", "foo": "bar"},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, stats.collectTags(test.tags))
	}
}
