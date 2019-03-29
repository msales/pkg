package cache

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisOptionsFunc represents an configuration function for Redis.
type RedisOptionsFunc func(*redis.UniversalOptions)

// WithPoolSize configures the Redis pool size.
func WithPoolSize(size int) RedisOptionsFunc {
	return func(o *redis.UniversalOptions) {
		o.PoolSize = size
	}
}

// WithPoolTimeout configures the Redis pool timeout.
func WithPoolTimeout(timeout time.Duration) RedisOptionsFunc {
	return func(o *redis.UniversalOptions) {
		o.PoolTimeout = timeout
	}
}

// WithReadTimeout configures the Redis read timeout.
func WithReadTimeout(timeout time.Duration) RedisOptionsFunc {
	return func(o *redis.UniversalOptions) {
		o.ReadTimeout = timeout
	}
}

// WithWriteTimeout configures the Redis write timeout.
func WithWriteTimeout(timeout time.Duration) RedisOptionsFunc {
	return func(o *redis.UniversalOptions) {
		o.WriteTimeout = timeout
	}
}

type redisCache struct {
	client  redis.UniversalClient
	decoder decoder
}

// NewRedis create a new Redis cache instance.
func NewRedis(uri string, opts ...RedisOptionsFunc) (Cache, error) {
	o, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	uo := &redis.UniversalOptions{
		Addrs:         []string{o.Addr},
		RouteRandomly: true,
	}

	for _, opt := range opts {
		opt(uo)
	}

	c := redis.NewUniversalClient(uo)

	return &redisCache{
		client:  c,
		decoder: stringDecoder{},
	}, nil
}

func NewRedisUniversal(addrs []string, opts ...RedisOptionsFunc) (Cache, error) {
	uo := &redis.UniversalOptions{
		Addrs:         addrs,
		RouteRandomly: true,
	}

	for _, opt := range opts {
		opt(uo)
	}

	c := redis.NewUniversalClient(uo)

	return &redisCache{
		client:  c,
		decoder: stringDecoder{},
	}, nil
}

// Get gets the item for the given key.
func (c redisCache) Get(key string) *Item {
	b, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		err = ErrCacheMiss
	}

	return &Item{
		decoder: c.decoder,
		value:   b,
		err:     err,
	}
}

// GetMulti gets the items for the given keys.
func (c redisCache) GetMulti(keys ...string) ([]*Item, error) {
	val, err := c.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	i := []*Item{}
	for _, v := range val {
		var err = ErrCacheMiss
		var b []byte
		if v != nil {
			b = []byte(v.(string))
			err = nil
		}

		i = append(i, &Item{
			decoder: c.decoder,
			value:   b,
			err:     err,
		})
	}

	return i, nil
}

// Set sets the item in the cache.
func (c redisCache) Set(key string, value interface{}, expire time.Duration) error {
	return c.client.Set(key, value, expire).Err()
}

// Add sets the item in the cache, but only if the key does not already exist.
func (c redisCache) Add(key string, value interface{}, expire time.Duration) error {
	if !c.client.SetNX(key, value, expire).Val() {
		return ErrNotStored
	}
	return nil
}

// Replace sets the item in the cache, but only if the key already exists.
func (c redisCache) Replace(key string, value interface{}, expire time.Duration) error {
	if !c.client.SetXX(key, value, expire).Val() {
		return ErrNotStored
	}
	return nil
}

// Delete deletes the item with the given key.
func (c redisCache) Delete(key string) error {
	return c.client.Del(key).Err()
}

// Inc increments a key by the value.
func (c redisCache) Inc(key string, value uint64) (int64, error) {
	return c.client.IncrBy(key, int64(value)).Result()
}

// Dec decrements a key by the value.
func (c redisCache) Dec(key string, value uint64) (int64, error) {
	return c.client.DecrBy(key, int64(value)).Result()
}
