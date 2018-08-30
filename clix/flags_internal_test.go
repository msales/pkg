package clix

import (
	"errors"
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_splitTags(t *testing.T) {
	tests := []struct {
		tags []string

		results []interface{}
		err     error
	}{
		{[]string{"a=b"}, []interface{}{"a", "b"}, nil},
		{[]string{"a=b", "c=d"}, []interface{}{"a", "b", "c", "d"}, nil},
		{[]string{"a"}, nil, errors.New("invalid tags string")},
		{[]string{"a=b", "c"}, nil, errors.New("invalid tags string")},
	}

	for _, tt := range tests {
		res, err := splitTags(tt.tags, "=")

		assert.Equal(t, res, tt.results)
		assert.Equal(t, err, tt.err)
	}
}
