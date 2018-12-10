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
	"gopkg.in/urfave/cli.v1"
)

func TestWithLogger(t *testing.T) {
	c, _ := newTestContext()
	ctx, err := clix.NewContext(c, clix.WithLogger(log.Null))
	assert.NoError(t, err)

	l, ok := log.FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, l, log.Null)
}

func TestWithStats(t *testing.T) {
	c, _ := newTestContext()
	ctx, err := clix.NewContext(c, clix.WithStats(stats.Null))
	assert.NoError(t, err)

	s, ok := stats.FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, s, stats.Null)
}

func TestNewContext(t *testing.T) {
	c, _ := newTestContext()

	ctx, err := clix.NewContext(c)

	assert.NoError(t, err)
	assert.IsType(t, &clix.Context{}, ctx)
}

func TestNewContext_LoggerError(t *testing.T) {
	c, flags := newTestContext()
	flags.String(clix.FlagLogFormat, "test", "")

	_, err := clix.NewContext(c)

	assert.Error(t, err)
}

func TestNewContext_StatsError(t *testing.T) {
	c, flags := newTestContext()
	flags.String(clix.FlagStatsDSN, "test://", "")

	_, err := clix.NewContext(c)

	assert.Error(t, err)
}

func TestContext_Close(t *testing.T) {
	c, _ := newTestContext()
	ctx, _ := clix.NewContext(c, clix.WithLogger(log.Null), clix.WithStats(stats.Null))

	err := ctx.Close()

	assert.NoError(t, err)
}

func TestContext_CloseErrors(t *testing.T) {
	tests := []struct {
		name     string
		logErr   error
		statsErr error
	}{
		{
			name:     "No Error",
			logErr:   nil,
			statsErr: nil,
		},
		{
			name:     "Logger Error",
			logErr:   errors.New("test"),
			statsErr: nil,
		},
		{
			name:     "Stats Error",
			logErr:   nil,
			statsErr: errors.New("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := new(MockLogger)
			l.On("Close").Return(tt.logErr)

			s := new(MockStats)
			s.On("Close").Return(tt.statsErr)

			c, _ := newTestContext()
			ctx, err := clix.NewContext(c, clix.WithLogger(l), clix.WithStats(s))
			assert.NoError(t, err)

			err = ctx.Close()

			if tt.logErr != nil || tt.statsErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newTestContext() (*cli.Context, *flag.FlagSet) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := cli.NewContext(cli.NewApp(), fs, nil)

	return c, fs
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, ctx ...interface{}) {}

func (m *MockLogger) Info(msg string, ctx ...interface{}) {}

func (m *MockLogger) Warn(msg string, ctx ...interface{}) {}

func (m *MockLogger) Error(msg string, ctx ...interface{}) {}

func (m *MockLogger) Crit(msg string, ctx ...interface{}) {}

func (m *MockLogger) Close() error {
	args := m.Called()
	return args.Error(0)
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
