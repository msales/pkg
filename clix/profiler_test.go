package clix

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
	"time"
)

func TestRunProfiler_Enabled(t *testing.T) {
	ctx := new(CtxMock)
	ctx.On("Bool", FlagProfiler).Return(true)
	ctx.On("String", FlagProfilerPort).Return("62874")

	runProfiler(ctx)
	defer ProfilerServer.Shutdown(context.Background())

	time.Sleep(10 * time.Millisecond)

	conn, err := net.DialTimeout("tcp",":62874", time.Second)
	assert.NoError(t, err)

	if err == nil {
		conn.Close()
	}
}

func TestRunProfiler_Disabled(t *testing.T) {
	ctx := new(CtxMock)
	ctx.On("Bool", FlagProfiler).Return(false)

	runProfiler(ctx)

	_, err := net.DialTimeout("tcp",":62874", time.Second)

	assert.Error(t, err)
}

type CtxMock struct {
	mock.Mock
}

func (m *CtxMock) Bool(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *CtxMock) String(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *CtxMock) StringSlice(name string) []string {
	args := m.Called(name)
	return args.Get(0).([]string)
}
