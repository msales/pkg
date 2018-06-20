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

func TestWithStatsFunc(t *testing.T) {
	tests := []struct {
		ctx    context.Context
		expect Stats
	}{
		{context.Background(), Null},
	}

	for _, tt := range tests {
		withStats(tt.ctx, func(s Stats) error {
			assert.Equal(t, tt.expect, s)
			return nil
		})
	}
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

func BenchmarkCollectTags(b *testing.B) {
	tags := map[string]string{
		"test1": "test",
		"test2": "test",
		"test3": "test",
		"test4": "test",
		"test5": "test",

	}
	addedTags := map[string]string{
		"k1":  "v",
		"k2":  "v",
		"k3":  "v",
		"k4":  "v",
		"k5":  "v",
	}
	s := NewTaggedStats(Null, tags)

	for n := 0; n < b.N; n++ {
		s.collectTags(addedTags)
	}
}

func TestNullStats(t *testing.T) {
	s := Null

	assert.Nil(t, s.Inc("test", 1, 1.0, nil))
	assert.Nil(t, s.Dec("test", 1, 1.0, nil))
	assert.Nil(t, s.Gauge("test", 1.0, 1.0, nil))
	assert.Nil(t, s.Timing("test", 0, 1.0, nil))
}
