package cache

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestWithPoolSize(t *testing.T) {
	o := &redis.UniversalOptions{}

	WithPoolSize(12)(o)

	assert.Equal(t, 12, o.PoolSize)
}

func TestWithPoolTimeout(t *testing.T) {
	o := &redis.UniversalOptions{}

	WithPoolTimeout(time.Second)(o)

	assert.Equal(t, time.Second, o.PoolTimeout)
}

func TestWithReadTimeout(t *testing.T) {
	o := &redis.UniversalOptions{}

	WithReadTimeout(time.Second)(o)

	assert.Equal(t, time.Second, o.ReadTimeout)
}

func TestWithWriteTimeout(t *testing.T) {
	o := &redis.UniversalOptions{}

	WithWriteTimeout(time.Second)(o)

	assert.Equal(t, time.Second, o.WriteTimeout)
}

func TestNewRedis(t *testing.T) {
	v, err := NewRedis([]string{"redis://test"}, WithPoolSize(12))
	c := v.(*redisCache)

	client, isCluster := c.client.(*redis.ClusterClient)
	if isCluster {
		assert.Equal(t, 12, client.Options().PoolSize)
	} else {
		client, ok := c.client.(*redis.Client)
		assert.True(t, ok)
		assert.Equal(t, 12, client.Options().PoolSize)
	}

	assert.NoError(t, err)
}
