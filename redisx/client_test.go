package redisx_test

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/msales/pkg/v3/redisx"
	"github.com/stretchr/testify/assert"
)

func TestClusterScanIterator_Next_StandardClient(t *testing.T) {
	client := getClient()
	match := "test"

	client.Set(match, "", 0)

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	n := scanIterator.Next()
	assert.True(t, n)
}

func TestClusterScanIterator_WithError(t *testing.T) {
	client := getClusterClient()
	match := "test"

	_, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.Error(t, err)
}

func getClient() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	c := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return c
}

func getClusterClient() *redis.ClusterClient {
	s1, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	s2, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{s1.Addr(),s2.Addr()},
	})

	return c
}
