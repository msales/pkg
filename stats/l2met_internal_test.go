package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatL2metTags(t *testing.T) {
	tags := []interface{}{
		"test", "test",
		"foo", "bar",
		"test", "baz",
	}

	assert.Equal(t, "", formatL2metTags(nil))
	assert.Equal(t, "", formatL2metTags([]interface{}{}))

	got := formatL2metTags(tags)
	assert.NotContains(t, got, "test=test")
	assert.Contains(t, got, "test=baz")
	assert.Contains(t, got, "foo=bar")
}
