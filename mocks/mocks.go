package mocks

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/mock"
	"gopkg.in/urfave/cli.v1"

	"github.com/msales/pkg/v3/clix"
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
	cliArgs := os.Args[0:1]
	for k, v := range args {
		cliArgs = append(cliArgs, fmt.Sprintf("-%s=%s", k, v))
	}

	var cCtx *cli.Context
	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(c *cli.Context) { cCtx = c }
	err := app.Run(cliArgs)
	if err != nil {
		panic(err)
	}

	return cCtx
}
