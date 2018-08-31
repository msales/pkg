package clix_test

import (
	"errors"
	"testing"

	"github.com/msales/pkg/clix"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestFlags_Merge(t *testing.T) {
	f1 := cli.StringFlag{}
	f2 := cli.StringFlag{}
	flags := clix.Flags{f1}
	newFlags := clix.Flags{f2}

	merged := flags.Merge(newFlags)

	assert.IsType(t, clix.Flags{}, merged)
	assert.Len(t, merged, 2)
	assert.Contains(t, flags, f1)
	assert.Contains(t, flags, f2)
}

func Test_SplitTags(t *testing.T) {
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
		res, err := clix.SplitTags(tt.tags, "=")

		assert.Equal(t, res, tt.results)
		assert.Equal(t, err, tt.err)
	}
}
