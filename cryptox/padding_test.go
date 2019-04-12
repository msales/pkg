package cryptox_test

import (
	"testing"

	"github.com/msales/pkg/v3/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestPKCS7Pad(t *testing.T) {
	tests := []struct {
		name      string
		blocksize int
		in        []byte
		want      []byte
		wantErr   bool
	}{
		{
			name:      "input fully divisible by the blocksize",
			blocksize: 4,
			in:        []byte{'t', 'e', 's', 't'},
			want:      []byte{'t', 'e', 's', 't', 0x4, 0x4, 0x4, 0x4},
		},
		{
			name:      "input less than the blocksize",
			blocksize: 4,
			in:        []byte{'t', 'e'},
			want:      []byte{'t', 'e', 0x2, 0x2},
		},
		{
			name:      "input more than the blocksize",
			blocksize: 4,
			in:        []byte{'t', 'e', 's', 't', 's'},
			want:      []byte{'t', 'e', 's', 't', 's', 0x3, 0x3, 0x3},
		},
		{
			name:      "0 blocksize",
			blocksize: 0,
			in:        []byte{'t', 'e', 's', 't'},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "negative blocksize",
			blocksize: -1,
			in:        []byte{'t', 'e', 's', 't'},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty input",
			blocksize: 4,
			in:        []byte{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "nil input",
			blocksize: 4,
			in:        nil,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := cryptox.PKCS7Pad(tt.in, tt.blocksize)

			assert.Equal(t, tt.want, out)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPKCS7Unpad(t *testing.T) {
	tests := []struct {
		name      string
		blocksize int
		in        []byte
		want      []byte
		wantErr   bool
	}{
		{
			name:      "full block to unpad",
			blocksize: 4,
			in:        []byte{'t', 'e', 's', 't', 0x4, 0x4, 0x4, 0x4},
			want:      []byte{'t', 'e', 's', 't'},
		},
		{
			name:      "part of block to unpad",
			blocksize: 4,
			in:        []byte{'t', 'e', 0x2, 0x2},
			want:      []byte{'t', 'e'},
		},
		{
			name:      "0 blocksize",
			blocksize: 0,
			in:        []byte{'t', 'e', 's', 't'},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "negative blocksize",
			blocksize: -1,
			in:        []byte{'t', 'e', 's', 't'},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty input",
			blocksize: 4,
			in:        []byte{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty input",
			blocksize: 4,
			in:        nil,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "input not fully divisible by the blocksize",
			blocksize: 4,
			in:        []byte{'t', 'e', 's', 't', 's'},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "zero padding byte",
			blocksize: 5,
			in:        []byte{'t', 'e', 's', 't', 0x0},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "zero padding byte greater than input length",
			blocksize: 5,
			in:        []byte{'t', 'e', 's', 't', 0x6},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "non-identical padding bytes",
			blocksize: 4,
			in:        []byte{'t', 'e', 0x3, 0x4},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := cryptox.PKCS7Unpad(tt.in, tt.blocksize)

			assert.Equal(t, tt.want, out)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func BenchmarkPKCS7Pad(b *testing.B) {
	var data = []byte("12345678")

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = cryptox.PKCS7Pad(data, 4)
	}
}

func BenchmarkPKCS7Unpad(b *testing.B) {
	var data = []byte{'1', '2', '3', '4', '5', '6', '7', '8', 0x4, 0x4, 0x4, 0x4}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = cryptox.PKCS7Unpad(data, 4)
	}
}
