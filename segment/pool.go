package segment

import (
	"crypto/md5"
	"hash"
	"sync"
)

// hashPool represents a pool of hashes.
type hashPool struct {
	sync.Pool
}

// GetHash returns a hash from the pool or creates a new one if noe are available.
func (p hashPool) GetHash() hash.Hash {
	if v := p.Get(); v != nil {
		h := v.(hash.Hash)
		h.Reset()

		return h
	}

	return md5.New()
}
