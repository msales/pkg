package cache_test

import (
	"net"
	"testing"

	"github.com/msales/pkg/v4/cache"
	"github.com/stretchr/testify/assert"
)

var (
	testRedisServer = "localhost:6379"
	skipRedis       = false
)

func init() {
	c, err := net.Dial("tcp", testRedisServer)
	if err != nil {
		skipRedis = true
		return
	}
	c.Write([]byte("SELECT 1\r\n"))
	c.Write([]byte("FLUSHDB\r\n"))
	c.Close()
}

func TestRedisCache(t *testing.T) {
	if skipRedis {
		t.Skipf("skipping test; no running server at %s", testRedisServer)
	}

	c, err := cache.NewRedis("redis://" + testRedisServer + "/1")
	assert.NoError(t, err)

	runCacheTests(t, c)
}
