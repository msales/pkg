package cache

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
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
	assert.NoError(t, err)

	c := v.(*redisCache)
	client := c.client.(*redis.Client)
	assert.IsType(t, &redis.Client{}, client)
	assert.Equal(t, 12, client.Options().PoolSize)
}

func TestNewRedisUniversal_Cluster(t *testing.T) {
	v, err := NewRedisUniversal([]string{"test:8179", "test2:8179"}, WithPoolSize(12))
	assert.NoError(t, err)

	c := v.(*redisCache)
	clusterClient := c.client.(*redis.ClusterClient)
	assert.IsType(t, &redis.ClusterClient{}, clusterClient)
	assert.Equal(t, 12, clusterClient.Options().PoolSize)
}

func TestNewRedisUniversal_NonCluster(t *testing.T) {
	v, err := NewRedisUniversal([]string{"test:8179"}, WithPoolSize(12))
	assert.NoError(t, err)

	c := v.(*redisCache)
	client := c.client.(*redis.Client)
	assert.IsType(t, &redis.Client{}, client)
	assert.Equal(t, 12, client.Options().PoolSize)
}

func TestNewRedis_InvalidUri(t *testing.T) {
	_, err := NewRedis("test")
	assert.Error(t, err)
}
