package clix_test

import (
	"net"
	"testing"
	"time"

	"github.com/msales/pkg/v3/clix"
	"github.com/stretchr/testify/assert"
)

func TestRunProfiler_Enabled(t *testing.T) {
	c, fs := newTestContext()
	fs.Bool(clix.FlagProfiler, true, "doc")
	fs.String(clix.FlagProfilerPort, "62874", "doc")

	clix.RunProfiler(c)
	defer clix.StopProfiler()

	time.Sleep(10 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", ":62874", time.Second)
	assert.NoError(t, err)

	if err == nil {
		conn.Close()
	}
}

func TestRunProfiler_Disabled(t *testing.T) {
	c, _ := newTestContext()

	clix.RunProfiler(c)

	_, err := net.DialTimeout("tcp", ":62874", time.Second)

	assert.Error(t, err)
}
