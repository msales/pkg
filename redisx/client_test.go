package redisx_test

import (
	"errors"
	"testing"

	"github.com/msales/pkg/v3/redisx"
	"github.com/stretchr/testify/assert"
)

func TestClusterScanIterator_Next_StandardClient(t *testing.T) {
	client := &clientMock{}
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	n := scanIterator.Next()
	assert.True(t, n)
}

func TestClusterScanIterator_Next_ClusterClient(t *testing.T) {
	client := &clusterClientMock{}
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	n := scanIterator.Next()
	assert.True(t, n)
}

func TestClusterScanIterator_Err(t *testing.T) {
	client := &clusterClientMock{}
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	err = scanIterator.Err()
	assert.NoError(t, err)
}

func TestClusterScanIterator_Next_WithTwoClientsFirstIteratorEnded(t *testing.T) {
	client := &EndedClusterClientMock{
		Masters: []redisx.Client{
			&endedClientMock{},
			&clientMock{},
		},
	}
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	ok := scanIterator.Next()
	assert.True(t, ok)
}

func TestClusterScanIterator_Next_WithTwoClientsBothIteratorEnded(t *testing.T) {
	client := &EndedClusterClientMock{
		Masters: []redisx.Client{
			&endedClientMock{},
			&endedClientMock{},
		},
	}
	match := "test"

	scanIterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)

	ok := scanIterator.Next()
	assert.False(t, ok)
}

func TestClusterScanIterator_WithError(t *testing.T) {
	client := &erroredClusterClient{}
	match := "test"

	_, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.Error(t, err)
}

func TestClusterScanIterator_Val_OnCluster(t *testing.T) {
	client := &clusterClientMock{}
	match := "test"

	iterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)
	assert.Equal(t, "value", iterator.Val())
}

func TestClusterScanIterator_Val_OnClient(t *testing.T) {
	client := &clientMock{}
	match := "test"

	iterator, err := redisx.NewScanIterator(client, 0, match, 0)
	assert.NoError(t, err)
	assert.Equal(t, "value", iterator.Val())
}

type erroredClusterClient struct {
}

func (c *erroredClusterClient) ForEachMaster(fn func(client redisx.Client) error) error {
	return errors.New("error")
}

type closedScanCmdMock struct {
}

func (*closedScanCmdMock) Iterator() redisx.ScanIterator {
	return &EndedScanIteratorMock{}
}

type ScanCmdMock struct {
}

func (s *ScanCmdMock) Iterator() redisx.ScanIterator {
	return &ScanIteratorMock{}
}

type ScanIteratorMock struct {
}

func (s *ScanIteratorMock) Val() string {
	return "value"
}

func (s *ScanIteratorMock) Next() bool {
	return true
}

func (s *ScanIteratorMock) Err() error {
	return nil
}

type EndedScanIteratorMock struct {
}

func (s *EndedScanIteratorMock) Val() string {
	return "value"
}

func (s *EndedScanIteratorMock) Next() bool {
	return false
}

func (s *EndedScanIteratorMock) Err() error {
	return nil
}

type clientMock struct {
}

func (c *clientMock) Scan(cursor uint64, match string, count int64) redisx.ScanCmd {
	return &ScanCmdMock{}
}

type EndedClusterClientMock struct {
	Masters []redisx.Client
}

func (e *EndedClusterClientMock) ForEachMaster(fn func(client redisx.Client) error) error {
	var err error
	for _, master := range e.Masters {
		err = fn(master)
	}

	return err
}

type endedClientMock struct {
}

func (*endedClientMock) Scan(cursor uint64, match string, count int64) redisx.ScanCmd {
	return &closedScanCmdMock{}
}

type clusterClientMock struct {
}

func (c *clusterClientMock) ForEachMaster(fn func(client redisx.Client) error) error {
	return fn(&clientMock{})
}
