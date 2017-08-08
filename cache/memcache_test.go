package cache_test

import (
	"net"
	"testing"

	"github.com/msales/pkg/cache"
	"github.com/stretchr/testify/assert"
)

var (
	testMemcachedServer = "localhost:11211"
	skip                = false
)

func init() {
	c, err := net.Dial("tcp", testMemcachedServer)
	if err != nil {
		skip = true
		return
	}
	c.Write([]byte("flush_all\r\n"))
	c.Close()
}

func TestMemcacheCache(t *testing.T) {
	if skip {
		t.Skipf("skipping test; no running server at %s", testMemcachedServer)
	}

	c := cache.NewMemcache(testMemcachedServer)

	// Set
	err := c.Set("test", "foobar", 0)
	assert.NoError(t, err)

	// Get
	str, err := c.Get("test").String()
	assert.NoError(t, err)
	assert.Equal(t, "foobar", str)
	_, err = c.Get("_").String()
	assert.EqualError(t, err, "cache: miss")

	// Add
	err = c.Add("test1", "foobar", 0)
	assert.NoError(t, err)
	err = c.Add("test1", "foobat", 0)
	assert.Error(t, err)

	// Replace
	err = c.Replace("test1", "foobar", 0)
	assert.NoError(t, err)
	err = c.Replace("test2", "foobar", 0)
	assert.Error(t, err)

	// GetMulti
	v, err := c.GetMulti("test", "test1", "_")
	assert.NoError(t, err)
	assert.Len(t, v, 3)
	assert.EqualError(t, v[2].Err(), "cache: miss")

	// Delete
	err = c.Delete("test1")
	assert.NoError(t, err)
	_, err = c.Get("test1").String()
	assert.Error(t, err)

	// Inc
	err = c.Set("test2", 1, 0)
	assert.NoError(t, err)
	i, err := c.Inc("test2", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	// Dec
	err = c.Set("test2", 1, 0)
	assert.NoError(t, err)
	i, err = c.Dec("test2", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), i)
}
