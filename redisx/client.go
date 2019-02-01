package redisx

type ScanIterator interface {
	Val() string
	Next() bool
	Err() error
}

type ClusterClient interface {
	ForEachMaster(fn func(client Client) error) error
}

type Client interface {
	Scan(cursor uint64, match string, count int64) ScanCmd
}

type ScanCmd interface {
	Iterator() ScanIterator
}

func NewScanIterator(c interface{}, cursor uint64, match string, count int64) (ScanIterator, error) {
	_, isCluster := c.(ClusterClient)

	if !isCluster {
		return c.(Client).Scan(cursor, match, count).Iterator(), nil
	}

	iterators := make([]ScanIterator, 0)
	err := c.(ClusterClient).ForEachMaster(func(client Client) error {
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

type ClusterScanIterator struct {
	iterators []ScanIterator

	curr int
}

func (cs *ClusterScanIterator) Val() string {
	return cs.getCurrentIterator().Val()
}

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
