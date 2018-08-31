package clix_test

import (
	"context"
	"github.com/msales/pkg/clix"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestRunProfiler_Enabled(t *testing.T) {
	ctx := new(CtxMock)
	ctx.On("Bool", clix.FlagProfiler).Return(true)
	ctx.On("Int", clix.FlagProfilerPort).Return(62874)

	clix.RunProfiler(ctx)
	defer clix.ProfilerServer.Shutdown(context.Background())

	time.Sleep(10 * time.Millisecond)

	conn, err := net.DialTimeout("tcp",":62874", time.Second)
	assert.NoError(t, err)

	if err == nil {
		conn.Close()
	}
}

func TestRunProfiler_Disabled(t *testing.T) {
	ctx := new(CtxMock)
	ctx.On("Bool", clix.FlagProfiler).Return(false)

	clix.RunProfiler(ctx)

	_, err := net.DialTimeout("tcp",":62874", time.Second)

	assert.Error(t, err)
}
