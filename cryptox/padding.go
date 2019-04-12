package cryptox

import (
	"bytes"
	"errors"
)

var (
	// errInvalidBlockSize indicates hash blocksize <= 0.
	errInvalidBlockSize = errors.New("cryptox: invalid blocksize")
	// errInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	errInvalidPKCS7Data = errors.New("cryptox: invalid PKCS7 data (empty or not padded)")
	// errInvalidPKCS7Padding indicates PKCS7 unpad fail because of bad input.
	errInvalidPKCS7Padding = errors.New("cryptox: invalid padding on input")
)

// PKCS7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func PKCS7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, errInvalidPKCS7Data
	}

	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)

	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))

	return pb, nil
}

// PKCS7Unpad validates and unpads THE data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func PKCS7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, errInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, errInvalidPKCS7Padding
	}

	c := b[len(b)-1]
	n := int(c)

	if n == 0 || n > len(b) {
		return nil, errInvalidPKCS7Padding
	}

	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, errInvalidPKCS7Padding
		}
	}

	return b[:len(b)-n], nil
}
