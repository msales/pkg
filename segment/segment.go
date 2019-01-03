package segment

import "encoding/hex"

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
	// ID is an arbitrary, unique identifier of an identity.
	ID string
	// Tier represents the length of the fields set.
	Tier int
	// Fields is a list of fields that make up the identity.
	Fields []string
}

// NewIdentity returns a new Identity instance.
func NewIdentity(id string, fields []string) Identity {
	return Identity{
		ID:     id,
		Tier:   len(fields),
		Fields: fields,
	}
}

// Segment represents a set of values defined by the identity.
type Segment struct {
	// Identity is the identity of the segment.
	Identity Identity
	// Values contains the byte values of the extracted fields. The order reflects the order of identity fields.
	Values   [][]byte
	// Hash is the hashed representation of the segment. It is nil if segment is not marshaled.
	Hash     []byte
}

// NewSegment returns a new Segment instance.
func NewSegment(identity Identity) Segment {
	return Segment{
		Identity: identity,
		Values:   make([][]byte, identity.Tier),
	}
}

// String returns a string representation of the segment.
func (s *Segment) String() string {
	return hex.EncodeToString(s.Hash)
}
