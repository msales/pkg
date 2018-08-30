package clix_test

import (
	"testing"

	"github.com/msales/pkg/clix"
	"github.com/stretchr/testify/assert"

	"github.com/urfave/cli"
)

func TestFlags(t *testing.T) {
	flags := clix.Flags()

	assert.IsType(t, []cli.Flag{}, flags)
	assert.NotEmpty(t, flags)
}

func TestProfilerFlags(t *testing.T) {
	flags := clix.ProfilerFlags()

	assert.IsType(t, []cli.Flag{}, flags)
	assert.NotEmpty(t, flags)
}
