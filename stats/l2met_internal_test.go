package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatL2metTags(t *testing.T) {
	tags := map[string]string{
		"test": "test",
		"foo":  "bar",
	}

	assert.Equal(t, "", formatL2metTags(nil))
	assert.Equal(t, "", formatL2metTags(map[string]string{}))

	got := formatL2metTags(tags)
	assert.Contains(t, got, "test=test")
	assert.Contains(t, got, "foo=bar")
}
