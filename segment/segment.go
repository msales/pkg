package segment

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	separator = []byte{'\n'}
	pool      = hashPool{}
)

// Marshal returns a segment for the given identity and value.
func Marshal(v SelfAccessor, identity Identity) (Segment, error) {
	return MarshalWith(Wrap(v), v, identity)
}

// MarshalWith returns a segment for the given identity and value, using the specified accessor.
func MarshalWith(accessor FieldAccessor, v interface{}, identity Identity) (Segment, error) {
	h := pool.GetHash()
	seg := NewSegment(identity)

	for i, field := range identity.Fields {
		b, err := accessor.GetField(v, field)
		if err != nil {
			return Segment{}, err
		}

		seg.Values[i] = b
		h.Write(b)
		h.Write(separator)
	}

	seg.Hash = h.Sum(nil)
	pool.Put(h)

	return seg, nil
}

// Identity represents a segment identity.
type Identity struct {
	// ID is an arbitrary, unique identifier of an identity. If segment hashes are used as keys, ID should also be part of it.
	// This is important in avoiding conflicts in cases where 2 similar identities might yield the same resulting segment.
	// E.g. Consider the following 2 identities:
	//				[a, b], [a, c]
	// If the data we're segmenting contains, coincidentally, the same value for the fields "b" and "c",
	// we will end up with identical segments. In that case the ID becomes the only property
	// that allows us to tell those segments apart.
	ID string
	// Fields is a list of fields that make up the identity.
	Fields []string
}

// NewIdentity returns a new Identity instance.
func NewIdentity(id string, fields []string) Identity {
	return Identity{
		ID:     id,
		Fields: fields,
	}
}

// Len returns the length of the identity (the number of fields that make it).
func (i Identity) Len() int {
	return len(i.Fields)
}

// Segment represents a set of values defined by the identity.
type Segment struct {
	// Identity is the identity of the segment.
	Identity Identity
	// Values contains the byte values of the extracted fields. The order reflects the order of identity fields.
	Values [][]byte
	// Hash is the hashed representation of the segment. It is nil if segment is not marshaled.
	Hash []byte
}

// NewSegment returns a new Segment instance.
func NewSegment(identity Identity) Segment {
	return Segment{
		Identity: identity,
		Values:   make([][]byte, identity.Len()),
	}
}

// String returns a string representation of the segment.
func (s *Segment) String() string {
	return hex.EncodeToString(s.Hash)
}

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
