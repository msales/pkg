package redisx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/msales/pkg/v3/redisx"
)

func TestClusterScanIterator_Next_StandardClient(t *testing.T) {
	client := getClient()
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	n := scanIterator.Next(context.Background())
	assert.False(t, n)
}

func TestClusterScanIterator_Next_ClusterClient(t *testing.T) {
	client1 := getClient()
	client2 := getClient()

	client1.Set(context.Background(), "test1", 1, 0)
	client2.Set(context.Background(), "test2", 2, 0)
	client2.Set(context.Background(), "test3", 3, 0)

	client := &clusterClientMock{
		Masters: []*redis.Client{
			client1,
			client2,
		},
	}
	match := "test*"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	n := scanIterator.Next(context.Background())
	assert.True(t, n)

	n = scanIterator.Next(context.Background())
	assert.True(t, n)

	n = scanIterator.Next(context.Background())
	assert.True(t, n)

	n = scanIterator.Next(context.Background())
	assert.False(t, n)
}

func TestClusterScanIterator_Err(t *testing.T) {
	client := &clusterClientMock{
		Masters: []*redis.Client{
			getClient(),
			getClient(),
		},
	}
	match := "test*"
	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	assert.NoError(t, scanIterator.Err())
}

func TestClusterScanIterator_WithError(t *testing.T) {
	client := &erroredClientMock{}
	match := "test*"

	_, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.Error(t, err)
}

func TestClusterScanIterator_Val(t *testing.T) {
	client1 := getClient()
	client2 := getClient()

	client1.Set(context.Background(), "test1", 1, 0)
	client2.Set(context.Background(), "test2", 2, 0)
	client2.Set(context.Background(), "test3", 3, 0)

	client := &clusterClientMock{
		Masters: []*redis.Client{
			client1,
			client2,
		},
	}
	match := "test*"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	assert.True(t, scanIterator.Next(context.Background()))
	assert.Equal(t, "test1", scanIterator.Val())
	assert.True(t, scanIterator.Next(context.Background()))
	assert.Equal(t, "test2", scanIterator.Val())
	assert.True(t, scanIterator.Next(context.Background()))
	assert.Equal(t, "test3", scanIterator.Val())
	assert.False(t, scanIterator.Next(context.Background()))

}

func getClient() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}

type clusterClientMock struct {
	Masters []*redis.Client
}

func (e *clusterClientMock) ForEachMaster(fn func(client *redis.Client) error) error {
	var err error
	for _, master := range e.Masters {
		err = fn(master)
	}

	return err
}

type erroredClientMock struct {
}

func (e *erroredClientMock) ForEachMaster(fn func(client *redis.Client) error) error {
	return errors.New("test error")
}
