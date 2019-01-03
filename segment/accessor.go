package segment

import "errors"

// common errors returned by this package.
var (
	ErrUndefinedField = errors.New("segment: undefined field")
)

// FieldAccessor represents a service that is able to access fields of other objects and return their byte value.
type FieldAccessor interface {
	// GetField returns a byte value of a required field, extracted from v.
	GetField(v interface{}, field string) ([]byte, error)
}

// SelfAccessor represents a FieldAccessor that operates on an internal state.
type SelfAccessor interface {
	// GetField returns a byte value of a required field.
	// Is should be extracted from the internal state.
	GetField(field string) ([]byte, error)
}

// FieldAccessorFunc represents a field accessor function.
type FieldAccessorFunc func(v interface{}, field string) ([]byte, error)

// GetField returns a byte value of a required field.
func (fn FieldAccessorFunc) GetField(v interface{}, field string) ([]byte, error) {
	return fn(v, field)
}

// SelfAccessorFunc represents a self accessor function.
type SelfAccessorFunc func(field string) ([]byte, error)

// GetField returns a byte value of a required field.
func (fn SelfAccessorFunc) GetField(field string) ([]byte, error) {
	return fn(field)
}

// StaticAccessor represents a map of static values that can be marshaled into a segment.
type StaticAccessor map[string][]byte

// GetField returns a byte value of a required field.
func (a StaticAccessor) GetField(field string) ([]byte, error) {
	val, ok := a[field]
	if !ok {
		return nil, ErrUndefinedField
	}

	return val, nil
}

// Wrap wraps the SelfAccessor instance and returns it as a FieldAccessor.
//
// Warning! Wrapped FieldAccessor will always use the GetField method of the wrapped value,
// ignoring the accessor passed as a first argument to it's own GetField!
func Wrap(v SelfAccessor) FieldAccessorFunc {
	return func(_ interface{}, field string) ([]byte, error) {
		return v.GetField(field)
	}
}
