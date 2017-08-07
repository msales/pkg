package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/msales/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	m := new(MockCache)
	m.On("Get", "test").Return(&cache.Item{})
	ctx := cache.WithCache(context.Background(), m)

	cache.Get(ctx, "test")

	m.AssertExpectations(t)
}

func TestGetMulti(t *testing.T) {
	m := new(MockCache)
	m.On("GetMulti", []string{"test"}).Return([]*cache.Item{&cache.Item{}}, nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.GetMulti(ctx, []string{"test"})

	m.AssertExpectations(t)
}

func TestSet(t *testing.T) {
	m := new(MockCache)
	m.On("Set", "test", 1, 0*time.Second).Return(nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Set(ctx, "test", 1, 0)

	m.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	m := new(MockCache)
	m.On("Add", "test", 1, 0*time.Second).Return(nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Add(ctx, "test", 1, 0)

	m.AssertExpectations(t)
}

func TestReplace(t *testing.T) {
	m := new(MockCache)
	m.On("Replace", "test", 1, 0*time.Second).Return(nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Replace(ctx, "test", 1, 0)

	m.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	m := new(MockCache)
	m.On("Delete", "test").Return(nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Delete(ctx, "test")

	m.AssertExpectations(t)
}

func TestInc(t *testing.T) {
	m := new(MockCache)
	m.On("Inc", "test", uint64(1)).Return(int64(1), nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Inc(ctx, "test", 1)

	m.AssertExpectations(t)
}

func TestDec(t *testing.T) {
	m := new(MockCache)
	m.On("Dec", "test", uint64(1)).Return(int64(1), nil)
	ctx := cache.WithCache(context.Background(), m)

	cache.Dec(ctx, "test", 1)

	m.AssertExpectations(t)
}

func TestNullCache_Get(t *testing.T) {
	i := cache.Null.Get("test")
	v, err := i.Bytes()

	assert.NoError(t, err)
	assert.Equal(t, []byte{}, v)
}

func TestNullCache_GetInt64(t *testing.T) {
	i := cache.Null.Get("test")
	b, err := i.Int64()

	assert.NoError(t, err)
	assert.Equal(t, int64(0), b)
}

func TestNullCache_GetFloat64(t *testing.T) {
	i := cache.Null.Get("test")
	v, err := i.Float64()

	assert.NoError(t, err)
	assert.Equal(t, float64(0), v)
}

func TestNullCache_GetMulti(t *testing.T) {
	v, err := cache.Null.GetMulti([]string{"test"})

	assert.NoError(t, err)
	assert.Len(t, v, 0)
}

func TestNullCache_Set(t *testing.T) {
	assert.NoError(t, cache.Null.Set("test", 1, 0))
}

func TestNullCache_Add(t *testing.T) {
	assert.NoError(t, cache.Null.Add("test", 1, 0))
}

func TestNullCache_Replace(t *testing.T) {
	assert.NoError(t, cache.Null.Replace("test", 1, 0))
}

func TestNullCache_Delete(t *testing.T) {
	assert.NoError(t, cache.Null.Delete("test"))
}

func TestNullCache_Inc(t *testing.T) {
	v, err := cache.Null.Inc("test", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)
}

func TestNullCache_Dec(t *testing.T) {
	v, err := cache.Null.Dec("test", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)
}

type MockCache struct {
	mock.Mock
}

func (c *MockCache) Get(key string) *cache.Item {
	args := c.Called(key)
	return args.Get(0).(*cache.Item)
}

func (c *MockCache) GetMulti(keys []string) ([]*cache.Item, error) {
	args := c.Called(keys)
	return args.Get(0).([]*cache.Item), args.Error(1)
}

func (c *MockCache) Set(key string, value interface{}, expire time.Duration) error {
	args := c.Called(key, value, expire)
	return args.Error(0)
}

func (c *MockCache) Add(key string, value interface{}, expire time.Duration) error {
	args := c.Called(key, value, expire)
	return args.Error(0)
}

func (c *MockCache) Replace(key string, value interface{}, expire time.Duration) error {
	args := c.Called(key, value, expire)
	return args.Error(0)
}

func (c *MockCache) Delete(key string) error {
	args := c.Called(key)
	return args.Error(0)
}

func (c *MockCache) Inc(key string, value uint64) (int64, error) {
	args := c.Called(key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (c *MockCache) Dec(key string, value uint64) (int64, error) {
	args := c.Called(key, value)
	return args.Get(0).(int64), args.Error(1)
}
