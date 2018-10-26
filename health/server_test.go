package health_test

import (
	"net"
	"testing"
	"time"

	"github.com/msales/pkg/health"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	r := &testReporter{}

	go health.StartServer(":62874", r)
	defer health.StopServer()

	conn, err := net.DialTimeout("tcp", ":62874", time.Second)
	assert.NoError(t, err)
	if err == nil {
		conn.Close()
	}
}
