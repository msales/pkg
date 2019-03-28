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
	v, err := NewRedis("redis://test", WithPoolSize(12))
	c := v.(*redisCache)

	assert.NoError(t, err)
	clusterClient, isCluster := c.client.(*redis.ClusterClient)
	if isCluster {
		assert.Equal(t, 12, clusterClient.Options().PoolSize)
	} else {
		client, ok := c.client.(*redis.Client)
		assert.True(t, ok)
		assert.Equal(t, 12, client.Options().PoolSize)
	}
}

func TestNewRedisUniversal(t *testing.T) {
	v, err := NewRedisUniversal([]string{"redis://test", "redis://test2"}, WithPoolSize(12))
	c := v.(*redisCache)

	assert.NoError(t, err)
	clusterClient, isCluster := c.client.(*redis.ClusterClient)
	if isCluster {
		assert.Equal(t, 12, clusterClient.Options().PoolSize)
	} else {
		client, ok := c.client.(*redis.Client)
		assert.True(t, ok)
		assert.Equal(t, 12, client.Options().PoolSize)
	}
}

func TestNewRedis_InvalidUri(t *testing.T) {
	_, err := NewRedis("test")
	assert.Error(t, err)
}
