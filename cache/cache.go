package cache

import (
	"context"
	"errors"
	"time"
)

type key int

const (
	ctxKey key = iota
)

var (
	// ErrCacheMiss means that a Get failed because the item wasn't present.
	ErrCacheMiss = errors.New("cache: miss")

	// ErrNotStored means the conditional write (Add or Replace) failed because
	// the condition was not met.
	ErrNotStored = errors.New("cache: not stored")

	// Null is the null Cache instance.
	Null = &nullCache{}
)

// Cache represents a cache instance.
type Cache interface {
	// Get gets the item for the given key.
	Get(key string) *Item

	// GetMulti gets the items for the given keys.
	GetMulti(keys ...string) ([]*Item, error)

	// Set sets the item in the cache.
	Set(key string, value interface{}, expire time.Duration) error

	// Add sets the item in the cache, but only if the key does not already exist.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets the item in the cache, but only if the key already exists.
	Replace(key string, value interface{}, expire time.Duration) error

	// Delete deletes the item with the given key.
	Delete(key string) error

	// Inc increments a key by the value.
	Inc(key string, value uint64) (int64, error)

	// Dec decrements a key by the value.
	Dec(key string, value uint64) (int64, error)
}

// WithCache sets Cache in the context.
func WithCache(ctx context.Context, cache Cache) context.Context {
	return context.WithValue(ctx, ctxKey, cache)
}

// FromContext returns the instance of Cache in the context.
func FromContext(ctx context.Context) (Cache, bool) {
	cache, ok := ctx.Value(ctxKey).(Cache)
	return cache, ok
}

// Get gets the item for the given key.
func Get(ctx context.Context, key string) *Item {
	c := getCache(ctx)
	return c.Get(key)
}

// GetMulti gets the items for the given keys.
func GetMulti(ctx context.Context, keys ...string) ([]*Item, error) {
	c := getCache(ctx)
	return c.GetMulti(keys...)
}

// Set sets the item in the cache.
func Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	c := getCache(ctx)
	return c.Set(key, value, expire)
}

// Add sets the item in the cache, but only if the key does not already exist.
func Add(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	c := getCache(ctx)
	return c.Add(key, value, expire)
}

// Replace sets the item in the cache, but only if the key already exists.
func Replace(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	c := getCache(ctx)
	return c.Replace(key, value, expire)
}

// Delete deletes the item with the given key.
func Delete(ctx context.Context, key string) error {
	c := getCache(ctx)
	return c.Delete(key)
}

// Inc increments a key by the value.
func Inc(ctx context.Context, key string, value uint64) (int64, error) {
	c := getCache(ctx)
	return c.Inc(key, value)
}

// Dec decrements a key by the value.
func Dec(ctx context.Context, key string, value uint64) (int64, error) {
	c := getCache(ctx)
	return c.Dec(key, value)
}

func getCache(ctx context.Context) Cache {
	if c, ok := FromContext(ctx); ok {
		return c
	}
	return Null
}

type nullDecoder struct{}

func (d nullDecoder) Bool(v []byte) (bool, error) {
	return false, nil
}

func (d nullDecoder) Int64(v []byte) (int64, error) {
	return 0, nil
}

func (d nullDecoder) Uint64(v []byte) (uint64, error) {
	return 0, nil
}

func (d nullDecoder) Float64(v []byte) (float64, error) {
	return 0, nil
}

type nullCache struct{}

// Get gets the item for the given key.
func (c nullCache) Get(key string) *Item {
	return &Item{decoder: nullDecoder{}, value: []byte{}}
}

// GetMulti gets the items for the given keys.
func (c nullCache) GetMulti(keys ...string) ([]*Item, error) {
	return []*Item{}, nil
}

// Set sets the item in the cache.
func (c nullCache) Set(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Add sets the item in the cache, but only if the key does not already exist.
func (c nullCache) Add(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Replace sets the item in the cache, but only if the key already exists.
func (c nullCache) Replace(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Delete deletes the item with the given key.
func (c nullCache) Delete(key string) error {
	return nil
}

// Inc increments a key by the value.
func (c nullCache) Inc(key string, value uint64) (int64, error) {
	return 0, nil
}

// Dec decrements a key by the value.
func (c nullCache) Dec(key string, value uint64) (int64, error) {
	return 0, nil
}
