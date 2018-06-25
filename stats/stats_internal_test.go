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

func TestNullStats(t *testing.T) {
	s := Null

	assert.Nil(t, s.Inc("test", 1, 1.0))
	assert.Nil(t, s.Dec("test", 1, 1.0))
	assert.Nil(t, s.Gauge("test", 1.0, 1.0))
	assert.Nil(t, s.Timing("test", 0, 1.0))

	assert.NoError(t, s.Close())
}

func TestMergeTags(t *testing.T) {
	tags := []interface{}{
		"test1", "foo",
		"test2", "bar",
	}

	tests := []struct {
		tags     []interface{}
		expected []interface{}
	}{
		{
			nil,
			[]interface{}{"test1", "foo", "test2", "bar"},
		},
		{
			[]interface{}{},
			[]interface{}{"test1", "foo", "test2", "bar"},
		},
		{
			[]interface{}{"foo", "bar", "baz", "bat"},
			[]interface{}{"test1", "foo", "test2", "bar", "foo", "bar", "baz", "bat"},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, mergeTags(tags, test.tags))
	}
}

func BenchmarkTaggedStats_MergeTags(b *testing.B) {
	tags := []interface{}{
		"test1", "test",
		"test2", "test",
		"test3", "test",
		"test4", "test",
		"test5", "test",
	}
	addedTags := []interface{}{
		"k1", "v",
		"k2", "v",
		"k3", "v",
		"k4", "v",
		"k5", "v",
	}

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		mergeTags(tags, addedTags)
	}
}

func TestDeduplicateTags(t *testing.T) {
	tests := []struct {
		tags     []interface{}
		expected []interface{}
	}{
		{
			[]interface{}{"test1", "foo", "test1", "bar"},
			[]interface{}{"test1", "bar"},
		},
		{
			[]interface{}{"test1", "foo", "test2", "bar"},
			[]interface{}{"test1", "foo", "test2", "bar"},
		},
		{
			[]interface{}{"test1", "foo", "test2", "bar", "test1", "baz"},
			[]interface{}{"test1", "baz", "test2", "bar"},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, deduplicateTags(test.tags))
	}
}

func BenchmarkTaggedStats_DeduplicateTags(b *testing.B) {
	tags := []interface{}{
		"test1", "foo",
		"test2", "bar",
		"test1", "baz",
		"test3", "test",
		"test4", "test",
		"test5", "test",
	}

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		deduplicateTags(tags)
	}
}
