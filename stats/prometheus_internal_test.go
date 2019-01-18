package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrometheus_formatFQN(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "test",
			want: "test",
		},
		{
			in:   "test",
			want: "test",
		},
		{
			in:   "test-test2",
			want: "test_test2",
		},
		{
			in:   "test.test.asdf.test",
			want: "test_test_asdf_test",
		},
		{
			in:   "test.test.asdf-test",
			want: "test_test_asdf_test",
		},
		{
			in:   "test.test.asdf.test.test.test.asdf.test.test.test.asdf.test.test.test.asdf..test",
			want: "test_test_asdf_test_test_test_asdf_test_test_test_asdf_test_test_test_asdf__test",
		},
		{
			in:   "test-test.asdf-test.test-test.asdf-test.test-test.asdf-test.test-test.asdf--test",
			want: "test_test_asdf_test_test_test_asdf_test_test_test_asdf_test_test_test_asdf__test",
		},
	}
	prometheus := NewPrometheus("test.test")

	for _, tt := range tests {
		res := prometheus.formatFQN(tt.in)
		assert.Equal(t, tt.want, res)
	}
}

func BenchmarkPrometheus_FormatFQN(b *testing.B) {
	testData := "test-test.asdf-test.test-test.asdf-test.test-test.asdf-test.test-test.asdf--test"

	prometheus := NewPrometheus("test")
	for i := 0; i < b.N; i++ {
		res := prometheus.formatFQN(testData)

		assert.NotNil(b, res)
	}
}
