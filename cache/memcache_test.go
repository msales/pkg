package cache_test

import (
	"net"
	"testing"

	"github.com/msales/pkg/v2/cache"
)

var (
	testMemcachedServer = "localhost:11211"
	skipMemcache        = false
)

func init() {
	c, err := net.Dial("tcp", testMemcachedServer)
	if err != nil {
		skipMemcache = true
		return
	}
	c.Write([]byte("flush_all\r\n"))
	c.Close()
}

func TestMemcacheCache(t *testing.T) {
	if skipMemcache {
		t.Skipf("skipping test; no running server at %s", testMemcachedServer)
	}

	c := cache.NewMemcache(testMemcachedServer)
	runCacheTests(t, c)
}
