package mocks

import (
	"fmt"

	"github.com/msales/pkg/v4/clix"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v2"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Debug(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}

func (m *Logger) Info(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}

func (m *Logger) Error(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}

// InitContext initializes clix context to be passed to existing application factories.
func InitContext(args map[string]string, flags []cli.Flag) *clix.Context {
	cCtx := InitCliContext(args, flags)

	ctx, err := clix.NewContext(cCtx)
	if err != nil {
		panic(err)
	}

	return ctx
}

func InitCliContext(args map[string]string, flags []cli.Flag) *cli.Context {
	cliArgs := []string{"test"}
	for k, v := range args {
		cliArgs = append(cliArgs, fmt.Sprintf("-%s=%s", k, v))
	}

	var cCtx *cli.Context
	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(c *cli.Context) error { cCtx = c; return nil }
	err := app.Run(cliArgs)
	if err != nil {
		panic(err)
	}

	return cCtx
}
