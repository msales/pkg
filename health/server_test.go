package health_test

import (
	"testing"
	"time"

	"github.com/msales/pkg/v4/health"
	"github.com/msales/pkg/v4/httpx"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	r := &testReporter{}

	go health.StartServer("127.0.0.1:8080", r)
	defer health.StopServer()

	time.Sleep(time.Millisecond)

	resp, err := httpx.Get("http://127.0.0.1:8080/health")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
