package mocks_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"

	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/mocks"
)

func TestInitContext(t *testing.T) {
	args := map[string]string{
		"foo": "bar",
		"faz": "3",
	}

	flags := clix.Flags{
		cli.StringFlag{
			Name:   "foo",
			EnvVar: "ENV_FOO",
		},
		cli.IntFlag{
			Name:   "faz",
			EnvVar: "ENV_FAZ",
		},
	}
	ctx := mocks.InitContext(args, flags)
	assert.Implements(t, (*context.Context)(nil), ctx)
	assert.IsType(t, &clix.Context{}, ctx)
}

func TestInitContext_PanicsBadFlag(t *testing.T) {
	args := map[string]string{
		"foo": "bar",
		"faz": "asdf",
	}

	flags := clix.Flags{
		cli.StringFlag{
			Name:   "foo",
			EnvVar: "ENV_FOO",
		},
		cli.IntFlag{
			Name:   "faz",
			EnvVar: "ENV_FAZ",
		},
	}
	assert.Panics(t, func() {
		_ = mocks.InitContext(args, flags)
	})
}

func TestInitContext_PanicsWrongLogLevel(t *testing.T) {
	args := map[string]string{
		"foo":             "bar",
		"faz":             "4",
		clix.FlagLogLevel: "not existing",
	}

	flags := clix.Flags{
		cli.StringFlag{
			Name:   "foo",
			EnvVar: "ENV_FOO",
		},
		cli.IntFlag{
			Name:   "faz",
			EnvVar: "ENV_FAZ",
		},
	}.Merge(clix.CommonFlags)
	assert.Panics(t, func() {
		_ = mocks.InitContext(args, flags)
	})
}
