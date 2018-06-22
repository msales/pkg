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

func TestFormatL2metTags_Tags(t *testing.T) {
	tags := Tags{}.
		With("test", "test").
		With("foo", "bar").
		With("test", "test")

	got := formatL2metTags([]interface{}{tags})
	assert.Contains(t, got, "test=test")
	assert.Contains(t, got, "foo=bar")
}

func TestFormatL2metTags_Uneven(t *testing.T) {
	tags := []interface{}{
		"test", "test",
		"foo",
	}

	defer func() {
		r := recover(); if r != nil {
			assert.NotNil(t, r)
		}
	}()

	formatL2metTags(tags)

	assert.Fail(t, "the test should have panicked on an uneven number of tags")
}
