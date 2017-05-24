package stats

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTaggedStats_CollectTags(t *testing.T) {
	stats := &TaggedStats{
		tags: map[string]string{
			"test1": "foo",
			"test2": "bar",
		},
	}

	tests := []struct {
		tags map[string]string
		expected map[string]string
	} {
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
