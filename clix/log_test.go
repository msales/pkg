package clix_test

import (
	"testing"

	"github.com/msales/pkg/v2/clix"
	"github.com/msales/pkg/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level  string
		format string
		tags   []string

		shouldErr bool
	}{
		{"info", "json", []string{}, false},
		{"", "json", []string{}, true},
		{"invalid", "json", []string{}, true},
		{"info", "", []string{}, true},
		{"info", "invalid", []string{}, true},
		{"info", "json", []string{"single"}, true},
	}

	for _, tt := range tests {
		ctx := new(CtxMock)
		ctx.On("String", clix.FlagLogLevel).Return(tt.level)
		ctx.On("String", clix.FlagLogFormat).Return(tt.format)
		ctx.On("StringSlice", clix.FlagLogTags).Return(tt.tags)

		l, err := clix.NewLogger(ctx)

		if tt.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Implements(t, (*log.Logger)(nil), l)
		}
	}
}

type CtxMock struct {
	mock.Mock
}

func (m *CtxMock) Bool(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *CtxMock) String(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *CtxMock) StringSlice(name string) []string {
	args := m.Called(name)
	return args.Get(0).([]string)
}
