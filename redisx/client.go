package redisx

import "github.com/go-redis/redis"

// ScanIterator represents a generic redis scan iterator that works on both
// redis Client and ClusterClient
type ScanIterator interface {
	Val() string
	Next() bool
	Err() error
}

// ScanCmd represents redis ScanCmd
type ScanCmd interface {
	Iterator() ScanIterator
}

type ClusterClient interface {
	ForEachMaster(fn func(client *redis.Client) error) error
}

// NewScanIterator returns a scan operator regarding redis client type
func NewScanIterator(c interface{}, cursor uint64, match string, count int64) (ScanIterator, error) {
	_, isCluster := c.(*redis.ClusterClient)

	if !isCluster {
		return c.(redis.Cmdable).Scan(cursor, match, count).Iterator(), nil
	}

	iterators := make([]ScanIterator, 0)
	err := c.(ClusterClient).ForEachMaster(func(client *redis.Client) error {
		iterators = append(iterators, client.Scan(cursor, match, count).Iterator())

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ClusterScanIterator{
		curr:      0,
		iterators: iterators,
	}, nil
}

// ClusterScanIterator represents redis cluster scan iterator
type ClusterScanIterator struct {
	iterators []ScanIterator

	curr int
}

// Val returns current value pointed by the iterator
func (cs *ClusterScanIterator) Val() string {
	return cs.getCurrentIterator().Val()
}

// Next returns true if there is at least one more value in iterator
func (cs *ClusterScanIterator) Next() bool {
	i := cs.getCurrentIterator()

	if i.Next() {
		return true
	}

	for cs.nextIterator() {
		i = cs.getCurrentIterator()

		if i.Next() {
			return true
		}
	}

	return false
}

// Err returns an error for iterator
func (cs *ClusterScanIterator) Err() error {
	return cs.getCurrentIterator().Err()
}

func (cs *ClusterScanIterator) getCurrentIterator() ScanIterator {
	return cs.iterators[0]
}

func (cs *ClusterScanIterator) nextIterator() bool {
	cs.iterators = cs.iterators[1:]

	return len(cs.iterators) > 0
}
