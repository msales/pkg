package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcacheOptionsFunc represents an configuration function for Memcache.
type MemcacheOptionsFunc func(*memcache.Client)

// WithIdleConns configures the Memcache max idle connections.
func WithIdleConns(size int) MemcacheOptionsFunc {
	return func(c *memcache.Client) {
		c.MaxIdleConns = size
	}
}

// WithTimeout configures the Memcache read and write timeout.
func WithTimeout(timeout time.Duration) MemcacheOptionsFunc {
	return func(c *memcache.Client) {
		c.Timeout = timeout
	}
}

type memcacheCache struct {
	client *memcache.Client

	encoder func(v interface{}) ([]byte, error)
	decoder decoder
}

// NewMemcache create a new Memcache cache instance.
func NewMemcache(uri string, opts ...MemcacheOptionsFunc) Cache {
	c := memcache.New(uri)

	for _, opt := range opts {
		opt(c)
	}

	return &memcacheCache{
		client:  c,
		encoder: memcacheEncoder,
		decoder: stringDecoder{},
	}
}

// Get gets the item for the given key.
func (c memcacheCache) Get(key string) *Item {
	b := []byte(nil)
	v, err := c.client.Get(key)
	switch err {
	case memcache.ErrCacheMiss:
		err = ErrCacheMiss
	case nil:
		b = v.Value
	}

	return &Item{
		decoder: c.decoder,
		value:   b,
		err:     err,
	}
}

// GetMulti gets the items for the given keys.
func (c memcacheCache) GetMulti(keys ...string) ([]*Item, error) {
	val, err := c.client.GetMulti(keys)
	if err != nil {
		return nil, err
	}

	i := []*Item{}
	for _, k := range keys {
		var err error = ErrCacheMiss
		var b []byte
		if v, ok := val[k]; ok {
			b = v.Value
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
func (c memcacheCache) Set(key string, value interface{}, expire time.Duration) error {
	v, err := c.encoder(value)
	if err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: int32(expire.Seconds()),
	})
}

// Add sets the item in the cache, but only if the key does not already exist.
func (c memcacheCache) Add(key string, value interface{}, expire time.Duration) error {
	v, err := c.encoder(value)
	if err != nil {
		return err
	}

	err = c.client.Add(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: int32(expire.Seconds()),
	})
	if err == memcache.ErrNotStored {
		return ErrNotStored
	}
	return err
}

// Replace sets the item in the cache, but only if the key already exists.
func (c memcacheCache) Replace(key string, value interface{}, expire time.Duration) error {
	v, err := c.encoder(value)
	if err != nil {
		return err
	}

	err = c.client.Replace(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: int32(expire.Seconds()),
	})

	if err == memcache.ErrNotStored {
		return ErrNotStored
	}
	return err
}

// Delete deletes the item with the given key.
func (c memcacheCache) Delete(key string) error {
	return c.client.Delete(key)
}

// Inc increments a key by the value.
func (c memcacheCache) Inc(key string, value uint64) (int64, error) {
	v, err := c.client.Increment(key, value)
	return int64(v), err
}

// Dec decrements a key by the value.
func (c memcacheCache) Dec(key string, value uint64) (int64, error) {
	v, err := c.client.Decrement(key, value)
	return int64(v), err
}

func memcacheEncoder(v interface{}) ([]byte, error) {
	switch v.(type) {
	case bool:
		if v.(bool) {
			return []byte("1"), nil
		}
		return []byte("0"), nil
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf("%d", v)), nil
	case uint, uint8, uint16, uint32, uint64:
		return []byte(fmt.Sprintf("%d", v)), nil
	case float32, float64:
		return []byte(fmt.Sprintf("%f", v)), nil
	case string:
		return []byte(v.(string)), nil
	}

	return json.Marshal(v)
}
