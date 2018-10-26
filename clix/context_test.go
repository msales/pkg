package clix_test

import (
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWithLogger(t *testing.T) {
	ctx := &clix.Context{}

	fn := clix.WithLogger(log.Null)
	assert.IsType(t, clix.ContextFunc(nil), fn)

	fn(ctx)

	l, ok := log.FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, l, log.Null)
}

func TestWithStats(t *testing.T) {
	ctx := &clix.Context{}

	fn := clix.WithStats(stats.Null)
	assert.IsType(t, clix.ContextFunc(nil), fn)

	fn(ctx)

	s, ok := stats.FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, s, stats.Null)
}

func TestNewContext(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("", flag.ContinueOnError), nil)

	ctx, err := clix.NewContext(c)

	assert.IsType(t, &clix.Context{}, ctx)
	assert.NoError(t, err)
}

func TestContext_Close(t *testing.T) {
	tests := []struct {
		err error
	}{
		{nil},
		{errors.New("")},
	}

	for _, tt := range tests {
		s := new(MockStats)
		s.On("Close").Return(tt.err)

		c := cli.NewContext(nil, flag.NewFlagSet("", flag.ContinueOnError), nil)
		ctx, err := clix.NewContext(c, clix.WithLogger(log.Null), clix.WithStats(s))
		assert.NoError(t, err)

		err = ctx.Close()

		assert.Equal(t, err, tt.err)
	}
}

func newTestContext() (*cli.Context, *flag.FlagSet) {
	fs := flag.NewFlagSet("test", 0)
	c := cli.NewContext(cli.NewApp(), fs, nil)

	return c, fs
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	return nil
}

func (m *MockStats) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	return nil
}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	return nil
}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	return nil
}

func (m *MockStats) Close() error {
	args := m.Called()
	return args.Error(0)
}
