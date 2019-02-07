package health_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msales/pkg/v3/health"
	"github.com/msales/pkg/v3/httpx"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	r := &testReporter{}

	mux := health.NewMux(r)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	time.Sleep(time.Millisecond)

	resp, err := httpx.Get( srv.URL+"/health")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
