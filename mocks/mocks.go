package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Debug(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}

func (m *Logger) Info(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}

func (m *Logger) Error(msg string, ctx ...interface{}) {
	args := []interface{}{msg}
	args = append(args, ctx...)
	m.Called(args...)
}
