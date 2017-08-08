package cache

import (
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
)

func TestWithIdleConns(t *testing.T) {
	c := &memcache.Client{}

	WithIdleConns(12)(c)

	assert.Equal(t, 12, c.MaxIdleConns)
}

func TestWithTimeout(t *testing.T) {
	c := &memcache.Client{}

	WithTimeout(time.Second)(c)

	assert.Equal(t, time.Second, c.Timeout)
}

func TestNewMemcache(t *testing.T) {
	c := NewMemcache("test", WithIdleConns(12)).(*memcacheCache)

	assert.Equal(t, 12, c.client.MaxIdleConns)
}

func TestByteEncode(t *testing.T) {
	tests := []struct {
		v      interface{}
		expect []byte
	}{
		{true, []byte("1")},
		{false, []byte("0")},
		{int64(10), []byte("10")},
		{uint64(10), []byte("10")},
		{float64(10.34), []byte("10.340000")},
		{"foobar", []byte("foobar")},
		{struct{ A int }{1}, []byte(`{"A":1}`)},
		{[]string{"foo", "bar"}, []byte(`["foo","bar"]`)},
	}

	for _, tt := range tests {
		got, err := byteEncode(tt.v)
		assert.NoError(t, err)

		assert.Equal(t, tt.expect, got)
	}
}
