package cache

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheOptionsFunc func(*memcache.Client)

func WithIdleConns(size int) MemcacheOptionsFunc {
	return func(c *memcache.Client) {
		c.MaxIdleConns = size
	}
}

func WithTimeout(timeout time.Duration) MemcacheOptionsFunc {
	return func(c *memcache.Client) {
		c.Timeout = timeout
	}
}

type memcacheCache struct {
	client  *memcache.Client
	decoder decoder
}

func NewMemcache(uri string, opts ...MemcacheOptionsFunc) Cache {
	c := memcache.New(uri)

	for _, opt := range opts {
		opt(c)
	}

	return &memcacheCache{
		client:  c,
		decoder: byteDecoder{},
	}
}

// Get gets the item for the given key.
func (c memcacheCache) Get(key string) *Item {
	v, err := c.client.Get(key)
	return &Item{
		decoder: c.decoder,
		value:   v.Value,
		err:     err,
	}
}

// GetMulti gets the items for the given keys.
func (c memcacheCache) GetMulti(keys []string) ([]*Item, error) {
	val, err := c.client.GetMulti(keys)
	if err != nil {
		return nil, err
	}

	i := []*Item{}
	for _, v := range val {
		i = append(i, &Item{
			decoder: c.decoder,
			value:   v.Value,
		})
	}

	return i, nil
}

// Set sets the item in the cache.
func (c memcacheCache) Set(key string, value interface{}, expire time.Duration) error {
	v, err := byteEncode(value)
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
	v, err := byteEncode(value)
	if err != nil {
		return err
	}

	return c.client.Add(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: int32(expire.Seconds()),
	})
}

// Replace sets the item in the cache, but only if the key already exists.
func (c memcacheCache) Replace(key string, value interface{}, expire time.Duration) error {
	v, err := byteEncode(value)
	if err != nil {
		return err
	}

	return c.client.Replace(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: int32(expire.Seconds()),
	})
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

type byteDecoder struct{}

func (d byteDecoder) Bool(b []byte) (bool, error) {
	var v bool
	err := byteDecode(b, v)
	return v, err
}

func (d byteDecoder) Int64(b []byte) (int64, error) {
	var v int64
	err := byteDecode(b, v)
	return v, err
}

func (d byteDecoder) Uint64(b []byte) (uint64, error) {
	var v uint64
	err := byteDecode(b, v)
	return v, err
}

func (d byteDecoder) Float64(b []byte) (float64, error) {
	var v float64
	err := byteDecode(b, v)
	return v, err
}

func byteEncode(v interface{}) ([]byte, error) {
	switch v.(type) {
	case int:
		v = int64(v.(int))
	case int8:
		v = int64(v.(int8))
	case int16:
		v = int64(v.(int16))
	case int32:
		v = int64(v.(int32))
	case uint:
		v = uint64(v.(uint))
	case uint8:
		v = uint64(v.(uint8))
	case uint16:
		v = uint64(v.(uint16))
	case uint32:
		v = uint64(v.(uint32))
	case float32:
		v = float64(v.(float32))
	}

	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func byteDecode(b []byte, v interface{}) error {
	err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, v)
	return err
}
