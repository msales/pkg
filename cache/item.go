package cache

import "strconv"

type decoder interface {
	Bool([]byte) (bool, error)
	Int64([]byte) (int64, error)
	Uint64([]byte) (uint64, error)
	Float64([]byte) (float64, error)
}

type stringDecoder struct{}

func (d stringDecoder) Bool(v []byte) (bool, error) {
	return string(v) == "1", nil
}

func (d stringDecoder) Int64(v []byte) (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func (d stringDecoder) Uint64(v []byte) (uint64, error) {
	return strconv.ParseUint(string(v), 10, 64)
}

func (d stringDecoder) Float64(v []byte) (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

// Item represents an item to be returned or stored in the cache
type Item struct {
	decoder decoder
	value   []byte
	err     error
}

// Bool gets the cache items value as a bool, or and error.
func (i Item) Bool() (bool, error) {
	if i.err != nil {
		return false, i.err
	}

	return i.decoder.Bool(i.value)
}

// Bytes gets the cache items value as bytes.
func (i Item) Bytes() ([]byte, error) {
	return i.value, i.err
}

// Bytes gets the cache items value as a string.
func (i Item) String() (string, error) {
	if i.err != nil {
		return "", i.err
	}

	return string(i.value), nil
}

// Int64 gets the cache items value as an int64, or and error.
func (i Item) Int64() (int64, error) {
	if i.err != nil {
		return 0, i.err
	}

	return i.decoder.Int64(i.value)
}

// Uint64 gets the cache items value as a uint64, or and error.
func (i Item) Uint64() (uint64, error) {
	if i.err != nil {
		return 0, i.err
	}

	return i.decoder.Uint64(i.value)
}

// Float64 gets the cache items value as a float64, or and error.
func (i Item) Float64() (float64, error) {
	if i.err != nil {
		return 0, i.err
	}

	return i.decoder.Float64(i.value)
}

// Err returns the item error or nil.
func (i Item) Err() error {
	return i.err
}
