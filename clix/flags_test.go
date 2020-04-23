package clix_test

import (
	"testing"

	"github.com/msales/pkg/v4/clix"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestFlags_Merge(t *testing.T) {
	f1 := &cli.StringFlag{}
	f2 := &cli.StringFlag{}
	f3 := &cli.StringFlag{}
	flags1 := clix.Flags{f1}
	flags2 := clix.Flags{f2}
	flags3 := clix.Flags{f3}

	merged := flags1.Merge(flags2, flags3)

	assert.IsType(t, clix.Flags{}, merged)
	assert.Len(t, merged, 3)
	assert.Contains(t, flags1, f1)
	assert.Contains(t, flags1, f2)
	assert.Contains(t, flags1, f3)
}
